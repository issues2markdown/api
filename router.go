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
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/issues2markdown/issues2markdown"
	"golang.org/x/oauth2"
)

// GETHomeOptions ...
func GETHomeOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization")
	w.WriteHeader(http.StatusOK)
}

// GETHome ...
func GETHome(w http.ResponseWriter, r *http.Request) {
	// Github Token
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) < 1 {
		log.Printf("ERROR: An Authorization header is required\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	githubToken := strings.Split(authorizationHeader, " ")[1]
	if githubToken == "" {
		log.Printf("ERROR: A valid Github Token is required\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := context.Background()

	// create github provider
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	issuesProvider := github.NewClient(tc)

	// create api client
	i2md, err := issues2markdown.NewIssuesToMarkdown(issuesProvider)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}

	log.Println("Querying data ...")
	qoptions := issues2markdown.NewQueryOptions(i2md.Username)

	// execute query
	issues, err := i2md.Query(qoptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rendering data ...")
	roptions := issues2markdown.NewRenderOptions()

	// render results
	result, err := i2md.Render(issues, roptions)
	if err != nil {
		log.Fatal(err)
	}

	// return response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/markdown; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		log.Fatal(err)
	}
}

// GETVersion ...
func GETVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	versionInfo := "issues2markdown API version info"
	_, err := w.Write([]byte(versionInfo))
	if err != nil {
		log.Fatal(err)
	}
}
