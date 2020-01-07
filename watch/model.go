package main

type meta struct {
	URL       string `json:"url"`
	RouteType int    `json:"route_type"`
	Remove    bool   `json:"remove"`
}
