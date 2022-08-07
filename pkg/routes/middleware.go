/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zinclabs/zinc/pkg/auth"
	"github.com/zinclabs/zinc/pkg/core"
)

func AuthMiddleware(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		if _, ok := auth.VerifyCredentials(user, password); ok {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Invalid credentials"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth": "Missing credentials"})
		return
	}
}

func ESMiddleware(c *gin.Context) {
	// Some es clients will check header("X-elastic-product") == "Elasticsearch".
	// If not, it will not work, and show "The client noticed that the server is not Elasticsearch and we do not support this unknown product."
	c.Header("X-elastic-product", "Elasticsearch")
}

func IndexAliasMiddleware(c *gin.Context) {
	target := ""
	ix := 0

	for i, entry := range c.Params {
		if entry.Key == "target" {
			target = entry.Value
			ix = i
			break
		}
	}

	if target == "" {
		c.Next()
		return
	}

	indexList := core.ZINC_INDEX_LIST.List()
	newTarget := ""

	// find all index that match this alias and add them to the newTarget
	for _, index := range indexList {
		if index.HasAlias(target) {
			newTarget += "," + index.GetName()
		}
	}

	if newTarget != "" {
		c.Params[ix].Value = newTarget // set target new value in the request context
	}
	c.Next()
}
