package terraform

import (
	"os"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func skipIntegration(t *testing.T) {
	if os.Getenv("C3_TEST_SKIP_INTEGRATION") != "" {
		t.Skip("Skipping Integration test...")
	}
}

func skipPlan(t *testing.T) {
	if os.Getenv("C3_TEST_SKIP_PLAN") != "" {
		t.Skip("Skipping TF-PLAN-BASED Test...")
	}
}

func skipAws(t *testing.T) {
	if os.Getenv("C3_TEST_SKIP_UNIT") != "" {
		t.Skip("Skipping tests requiring AWS credentials...")
	}
}

func WillActionByTagValue(plan tfjson.Plan, resource string, tagName string, tagValue string, changeAction string) bool {
	for _, changeItem := range plan.ResourceChanges {
		if changeItem.Type == resource {
			for _, changeStatus := range changeItem.Change.Actions {
				if changeAction == string(changeStatus) {
					if afterMap, ok := changeItem.Change.After.(map[string]interface{}); ok {
						for afterKey, afterVal := range afterMap {
							if afterKey == "tags" {
								if tagObject, ok := afterVal.(map[string]interface{}); ok {
									for awsTagName, awsTagValue := range tagObject {
										if tagName == awsTagName && tagValue == awsTagValue {
											return true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func WillCreateByName(plan tfjson.Plan, resource string, name string) bool {
	return WillActionByTagValue(plan, resource, "Name", name, "create")
}

func WillDestroyByName(plan tfjson.Plan, resource string, name string) bool {
	return WillActionByTagValue(plan, resource, "Name", name, "destroy")
}
