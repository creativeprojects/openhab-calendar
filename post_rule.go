package main

import "time"

func (l *Loader) PostRules(result string, date time.Time, rules []RuleConfiguration, postRules []PostRuleConfiguration) (string, error) {
	for _, postRule := range postRules {
		if postRule.When == nil {
			continue
		}
		if !MatchPostRule(result, *postRule.When) {
			continue
		}
		if postRule.Previous != nil {
			previous, err := l.GetResultFromRules(date.Add(-24*time.Hour), rules)
			if err != nil {
				return result, err
			}
			if !MatchPostRule(previous, *postRule.Previous) {
				continue
			}
		}
		if postRule.Next != nil {
			next, err := l.GetResultFromRules(date.Add(24*time.Hour), rules)
			if err != nil {
				return result, err
			}
			if !MatchPostRule(next, *postRule.Next) {
				continue
			}
		}
		// that's a match!
		return postRule.Result, nil
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
