// Copyright 2018 The issues2markdown Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/issues2markdown/issues2markdown"
	"golang.org/x/oauth2"
)

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Git Hub Token
		gitHubToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		s.Options.GitHubToken = gitHubToken

		if s.Options.GitHubToken == "" {
			http.Error(w, "A valid GitHub Token is required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()

		// create github client
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: s.Options.GitHubToken,
			},
		)
		tc := oauth2.NewClient(ctx, ts)
		issuesProvider := github.NewClient(tc)

		i2md, err := issues2markdown.NewIssuesToMarkdown(issuesProvider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Querying data ...")
		qoptions := issues2markdown.NewQueryOptions()
		qoptions.Organization = i2md.Username

		// execute query
		var args []string
		query := r.URL.Query().Get("q")
		if query != "" {
			args = append(args, query)
		}
		issues, err := i2md.Query(qoptions, strings.Join(args, " "))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Rendering data ...")
		roptions := issues2markdown.NewRenderOptions()
		result, err := i2md.Render(issues, roptions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/markdown; charset=UTF-8")
		fmt.Fprintf(w, result)
	}
}
