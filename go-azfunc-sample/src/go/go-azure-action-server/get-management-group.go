package main

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/resources/mgmt/managementgroups"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func getManagementGroupsClient() (managementgroups.Client, error) {
	managementGroupsClient := managementgroups.NewClient("", nil, nil, "")
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err == nil {
		managementGroupsClient.Authorizer = authorizer
		return managementGroupsClient, nil
	}

	return managementGroupsClient, err
}

func getManagementGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var invokeResponse InvokeResponse
	outputs := make(map[string]interface{})

	managementClient, err := getManagementGroupsClient()

	// Get scope from URL
	queryParams := r.URL.Query()
	managementGroupID := ""
	for k, v := range queryParams {
		if k == "managementGroupId" {
			managementGroupID = v[0]
		}
	}
	if len(managementGroupID) == 0 {
		http.Error(w, "No managementGroupId provided", http.StatusBadRequest)
		return
	}

	if err == nil {
		recurse := new(bool)
		*recurse = true
		group, err := managementClient.Get(ctx, managementGroupID, "children", recurse, "", "")

		var childrenArray []interface{}
		for _, v := range *group.Children {
			childrenArray = append(childrenArray, v)
		}
		outputs["Children"] = childrenArray

		if err == nil {
			outputs["Data"] = group

			invokeResponse = InvokeResponse{outputs, []string{"Management group get succeeded"}, http.StatusOK}
		} else {
			invokeResponse = InvokeResponse{outputs, []string{err.Error()}, http.StatusInternalServerError}
		}

		descentantResult, err := managementClient.GetDescendants(ctx, managementGroupID)
		if err == nil {
			var descendants []interface{}
			for _, v := range *descentantResult.Response().Value {
				descendants = append(descendants, v)
			}

			outputs["Descendants"] = descendants
		}
	} else {
		invokeResponse = InvokeResponse{outputs, []string{err.Error()}, http.StatusInternalServerError}
	}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
