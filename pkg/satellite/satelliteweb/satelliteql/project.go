// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package satelliteql

import (
	"github.com/graphql-go/graphql"
	"storj.io/storj/pkg/satellite"
)

const (
	projectType      = "project"
	projectInputType = "projectInput"

	fieldOwnerName   = "ownerName"
	fieldDescription = "description"
	// Indicates if user accepted Terms & Conditions during project creation
	// Used in input model
	fieldIsTermsAccepted = "isTermsAccepted"
)

// graphqlProject creates *graphql.Object type representation of satellite.ProjectInfo
func graphqlProject() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: projectType,
		Fields: graphql.Fields{
			fieldID: &graphql.Field{
				Type: graphql.String,
			},
			fieldOwnerName: &graphql.Field{
				Type: graphql.String,
			},
			fieldName: &graphql.Field{
				Type: graphql.String,
			},
			fieldDescription: &graphql.Field{
				Type: graphql.String,
			},
			fieldIsTermsAccepted: &graphql.Field{
				Type: graphql.Int,
			},
			fieldCreatedAt: &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}

// graphqlProjectInput creates graphql.InputObject type needed to create/update satellite.Project
func graphqlProjectInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: projectInputType,
		Fields: graphql.InputObjectConfigFieldMap{
			fieldID: &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			fieldName: &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			fieldDescription: &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			fieldIsTermsAccepted: &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
	})
}

// fromMapProjectInfo creates satellite.ProjectInfo from input args
func fromMapProjectInfo(args map[string]interface{}) (project satellite.ProjectInfo) {
	project.Name, _ = args[fieldName].(string)
	project.Description, _ = args[fieldDescription].(string)
	project.IsTermsAccepted, _ = args[fieldIsTermsAccepted].(bool)

	return
}
