// Author:: Jonathan Hartman (<j@p4nt5.com>)
//
// Copyright (C) 2014, Jonathan Hartman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cookbook_collection

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "github.com/RoboticCheese/supermarket-circular/cookbook"
)

type CookbookCollection struct {
  Cookbooks []cookbook.Cookbook
}

func (cc *CookbookCollection) NewFromUniverse(url string) (*CookbookCollection, error) {
  univ, err := cc.universe(url)
  if err != nil {
    return cc, err
  }
  cc, err = cc.cc_from_universe(univ)
  return cc, err
}

func (old_cc *CookbookCollection) Merge(new_cc CookbookCollection) (*CookbookCollection, error) {
  for _, cookbook := range new_cc.Cookbooks {
    if !old_cc.contains(cookbook) {
      old_cc.Cookbooks = append(old_cc.Cookbooks, cookbook)
    } else {
      // TODO: Make contains return an index to save the second loop
      for k, cb := range old_cc.Cookbooks {
        if cb.Name == cookbook.Name {
          old_cc.Cookbooks[k].Merge(cookbook)
        }
      }
    }
  }
  return old_cc, nil
}

func (cc *CookbookCollection) contains(cb cookbook.Cookbook) (bool) {
  for _, cur := range cc.Cookbooks {
    if cur.Name == cb.Name {
      return true
    }
  }
  return false
}

func (cc *CookbookCollection) cc_from_universe(universe []byte) (*CookbookCollection, error) {
  var i interface{}
  err := json.Unmarshal(universe, &i)
  if err != nil {
    return cc, err
  }

  m := i.(map[string]interface{})
  for name, attrs := range m {
    new_cb := cookbook.Cookbook{name, []string{}}
    versions := attrs.(map[string]interface{})
    for version, _ := range versions {
      new_cb.Versions = append(new_cb.Versions, version)
    }
    cc.Cookbooks = append(cc.Cookbooks, new_cb)
  }
  return cc, nil
}

func (cc *CookbookCollection) universe(url string) (body []byte, err error) {
  resp, err := http.Get(url)
  if err != nil {
    return
  }
  defer resp.Body.Close()

  body, err = ioutil.ReadAll(resp.Body)
  return
}
