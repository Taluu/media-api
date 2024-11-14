package adapters

import (
	mediaFake "github.com/Taluu/media-go/pkg/domain/media/adapters/media/fake"
	tagFake "github.com/Taluu/media-go/pkg/domain/media/adapters/tag/fake"
	uploaderFake "github.com/Taluu/media-go/pkg/domain/media/adapters/uploader/fake"
	uploaderFile "github.com/Taluu/media-go/pkg/domain/media/adapters/uploader/fake"
)

var (
	NewFakeMediaRepository = mediaFake.NewFake
	NewFakeTagRegistry     = tagFake.NewFake
	NewFakeUploader        = uploaderFake.NewUploader
	NewFileUploader        = uploaderFile.NewUploader
)
