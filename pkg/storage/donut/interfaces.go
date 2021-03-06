/*
 * Minimalist Object Storage, (C) 2015 Minio, Inc.
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

package donut

import (
	"io"

	"github.com/minio/minio/pkg/storage/donut/disk"
)

// Collection of Donut specification interfaces

// Donut is a collection of object storage and management interface
type Donut interface {
	ObjectStorage
	Management
}

// ObjectStorage is a donut object storage interface
type ObjectStorage interface {
	// Storage service operations
	GetBucketMetadata(bucket string) (map[string]string, error)
	SetBucketMetadata(bucket string, metadata map[string]string) error
	ListBuckets() (map[string]map[string]string, error)
	MakeBucket(bucket, acl string) error

	// Bucket operations
	ListObjects(bucket, prefix, marker, delim string, maxKeys int) (objects []string, prefixes []string, isTruncated bool, err error)

	// Object operations
	GetObject(bucket, object string) (io.ReadCloser, int64, error)
	GetObjectMetadata(bucket, object string) (map[string]string, error)
	PutObject(bucket, object, expectedMD5Sum string, reader io.ReadCloser, metadata map[string]string) (string, error)
}

// Management is a donut management system interface
type Management interface {
	Heal() error
	Rebalance() error
	Info() (map[string][]string, error)

	AttachNode(node Node) error
	DetachNode(node Node) error

	SaveConfig() error
	LoadConfig() error
}

// Node interface for node management
type Node interface {
	ListDisks() (map[int]disk.Disk, error)
	AttachDisk(disk disk.Disk, diskOrder int) error
	DetachDisk(diskOrder int) error

	GetNodeName() string
	SaveConfig() error
	LoadConfig() error
}
