package main

type Repo struct {
	Owner   string
	Project string
	Package string
}

type BuildListResp struct {
	Items []struct {
		ID            int    `json:"id"`
		State         string `json:"state"`
		SourcePackage struct {
			Name string `json:"name"`
		} `json:"source_package"`
	} `json:"items"`
}
