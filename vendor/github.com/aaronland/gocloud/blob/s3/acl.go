package s3

import (
	"fmt"

	tm_types "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	s3_types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// StringACLToObjectCannedACL resolves string values for S3 ACLs to their corresponding
// `github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types.ObjectCannedACL` instance.
func StringACLToTransferManagerObjectCannedACL(str_acl string) (tm_types.ObjectCannedACL, error) {

	switch str_acl {
	case "private":
		return tm_types.ObjectCannedACLPrivate, nil
	case "public-read":
		return tm_types.ObjectCannedACLPublicRead, nil
	case "public-read-write":
		return tm_types.ObjectCannedACLPublicReadWrite, nil
	case "authenticated-read":
		return tm_types.ObjectCannedACLAuthenticatedRead, nil
	case "aws-exec-read":
		return tm_types.ObjectCannedACLAwsExecRead, nil
	case "bucket-owner-read":
		return tm_types.ObjectCannedACLBucketOwnerRead, nil
	case "bucket-owner-full-control":
		return tm_types.ObjectCannedACLBucketOwnerFullControl, nil
	default:
		return "", fmt.Errorf("Unsupported ACL string")
	}

}

// StringACLToObjectCannedACL resolves a subset of the string values for S3 ACLs (those specific to objects) to
// their corresponding `github.com/aws/aws-sdk-go-v2/service/s3/types.ObjectCannedACL` instance.
func StringACLToObjectCannedACL(str_acl string) (s3_types.ObjectCannedACL, error) {

	switch str_acl {
	case "private":
		return s3_types.ObjectCannedACLPrivate, nil
	case "public-read":
		return s3_types.ObjectCannedACLPublicRead, nil
	case "public-read-write":
		return s3_types.ObjectCannedACLPublicReadWrite, nil
	case "authenticated-read":
		return s3_types.ObjectCannedACLAuthenticatedRead, nil
	default:
		return "", fmt.Errorf("Invalid or unsupported ACL")
	}

}
