package main

import "time"

func (l *Loader) PostRules(result Result, date time.Time, rules []RuleConfiguration, postRules []PostRuleConfiguration) (Result, error) {
	for _, postRule := range postRules {
		if postRule.When == nil {
			continue
		}
		if !MatchPostRule(result.Calendar, *postRule.When) {
			continue
		}
		if postRule.Previous != nil {
			previous, err := l.GetResultFromRules(date.Add(-24*time.Hour), rules)
			if err != nil {
				return result, err
			}
			if !MatchPostRule(previous.Calendar, *postRule.Previous) {
				continue
			}
		}
		if postRule.Next != nil {
			next, err := l.GetResultFromRules(date.Add(24*time.Hour), rules)
			if err != nil {
				return result, err
			}
			if !MatchPostRule(next.Calendar, *postRule.Next) {
				continue
			}
		}
		// that's a match!
		return Result{Calendar: postRule.Result}, nil
	}
	return result, nil
}

func MatchPostRule(value string, matcher PostRuleMatcher) bool {
	if matcher.Is != "" {
		if value != matcher.Is {
			return false
		}
	}
	if matcher.Not != "" {
		if value == matcher.Not {
			return false
		}
	}
	return true
}
