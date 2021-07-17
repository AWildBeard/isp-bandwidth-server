package main

type Error struct {
	Message string `json:"message"`
}

type UploadStatus struct {
	Ok bool `json:"ok"`
}