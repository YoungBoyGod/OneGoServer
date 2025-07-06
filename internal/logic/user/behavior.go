package user

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/os/gtime"

	"OneGfServer/internal/model/user"
)

// ===============================
// 用户行为分析相关业务逻辑
// ===============================

// AnalyzeUserBehavior 分析用户行为模式
func (s *sUser) AnalyzeUserBehavior(ctx context.Context, input *user.AnalyzeUserBehaviorInput) (*user.AnalyzeUserBehaviorOutput, error) {
	if len(input.ActivityData) == 0 {
		return &user.AnalyzeUserBehaviorOutput{
			Analysis: map[string]interface{}{
				"status":  "insufficient_data",
				"message": "活动数据不足",
			},
		}, nil
	}

	analysis := make(map[string]interface{})

	// 统计活动类型分布
	actionCounts := make(map[string]int)
	hourCounts := make(map[int]int)
	var totalSessions, totalDuration float64

	for _, activity := range input.ActivityData {
		// 统计操作类型
		if action, ok := activity["action_type"].(string); ok {
			actionCounts[action]++
		}

		// 统计活跃时间段
		if timestamp, ok := activity["timestamp"].(*gtime.Time); ok && timestamp != nil {
			hour := timestamp.Time.Hour()
			hourCounts[hour]++
		}

		// 统计会话时长
		if duration, ok := activity["session_duration"].(float64); ok {
			totalDuration += duration
			totalSessions++
		}
	}

	analysis["action_distribution"] = actionCounts
	analysis["active_hours"] = hourCounts

	if totalSessions > 0 {
		analysis["avg_session_duration"] = totalDuration / totalSessions
	}

	// 识别活跃时间段
	mostActiveHourInput := &user.FindMostActiveHourInput{
		HourCounts: hourCounts,
	}
	mostActiveHourOutput := s.FindMostActiveHour(ctx, mostActiveHourInput)
	analysis["most_active_hour"] = mostActiveHourOutput.MostActiveHour

	// 识别行为模式
	behaviorPatternInput := &user.IdentifyBehaviorPatternInput{
		ActionCounts: actionCounts,
		HourCounts:   hourCounts,
	}
	behaviorPatternOutput := s.IdentifyBehaviorPattern(ctx, behaviorPatternInput)
	analysis["behavior_pattern"] = behaviorPatternOutput.Pattern

	// 计算活跃度评分
	activityScoreInput := &user.CalculateActivityScoreInput{
		ActivityData: input.ActivityData,
	}
	activityScoreOutput := s.CalculateActivityScore(ctx, activityScoreInput)
	analysis["activity_score"] = activityScoreOutput.Score

	return &user.AnalyzeUserBehaviorOutput{
		Analysis: analysis,
	}, nil
}

// FindMostActiveHour 找到最活跃的时间段
func (s *sUser) FindMostActiveHour(ctx context.Context, input *user.FindMostActiveHourInput) (*user.FindMostActiveHourOutput, error) {
	maxCount := 0
	mostActiveHour := 0

	for hour, count := range input.HourCounts {
		if count > maxCount {
			maxCount = count
			mostActiveHour = hour
		}
	}

	return &user.FindMostActiveHourOutput{
		MostActiveHour: mostActiveHour,
	}, nil
}

// IdentifyBehaviorPattern 识别行为模式
func (s *sUser) IdentifyBehaviorPattern(ctx context.Context, input *user.IdentifyBehaviorPatternInput) (*user.IdentifyBehaviorPatternOutput, error) {
	// 根据活动类型判断用户类型
	readCount := input.ActionCounts["view"] + input.ActionCounts["read"]
	writeCount := input.ActionCounts["create"] + input.ActionCounts["update"] + input.ActionCounts["delete"]

	var pattern string
	if writeCount > readCount*2 {
		pattern = "power_user"
	} else if readCount > writeCount*5 {
		pattern = "viewer"
	} else {
		pattern = "normal_user"
	}

	return &user.IdentifyBehaviorPatternOutput{
		Pattern: pattern,
	}, nil
}

// CalculateActivityScore 计算活跃度评分
func (s *sUser) CalculateActivityScore(ctx context.Context, input *user.CalculateActivityScoreInput) (*user.CalculateActivityScoreOutput, error) {
	if len(input.ActivityData) == 0 {
		return &user.CalculateActivityScoreOutput{
			Score: 0,
		}, nil
	}

	score := 0.0

	// 基于活动频率
	activityPerDay := float64(len(input.ActivityData)) / 30.0 // 假设30天的数据
	if activityPerDay > 10 {
		score += 50
	} else {
		score += activityPerDay * 5
	}

	// 基于活动多样性
	actionTypes := make(map[string]bool)
	for _, activity := range input.ActivityData {
		if action, ok := activity["action_type"].(string); ok {
			actionTypes[action] = true
		}
	}
	score += float64(len(actionTypes)) * 5

	return &user.CalculateActivityScoreOutput{
		Score: math.Min(score, 100),
	}, nil
}
