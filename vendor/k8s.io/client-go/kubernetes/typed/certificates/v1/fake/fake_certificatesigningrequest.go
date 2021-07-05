/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	certificatesv1 "k8s.io/api/certificates/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCertificateSigningRequests implements CertificateSigningRequestInterface
type FakeCertificateSigningRequests struct {
	Fake *FakeCertificatesV1
}

var certificatesigningrequestsResource = schema.GroupVersionResource{Group: "certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"}

var certificatesigningrequestsKind = schema.GroupVersionKind{Group: "certificates.k8s.io", Version: "v1", Kind: "CertificateSigningRequest"}

// Get takes name of the certificateSigningRequest, and returns the corresponding certificateSigningRequest object, and an error if there is any.
func (c *FakeCertificateSigningRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *certificatesv1.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(certificatesigningrequestsResource, name), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}

// List takes label and field selectors, and returns the list of CertificateSigningRequests that match those selectors.
func (c *FakeCertificateSigningRequests) List(ctx context.Context, opts v1.ListOptions) (result *certificatesv1.CertificateSigningRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(certificatesigningrequestsResource, certificatesigningrequestsKind, opts), &certificatesv1.CertificateSigningRequestList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &certificatesv1.CertificateSigningRequestList{ListMeta: obj.(*certificatesv1.CertificateSigningRequestList).ListMeta}
	for _, item := range obj.(*certificatesv1.CertificateSigningRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificateSigningRequests.
func (c *FakeCertificateSigningRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(certificatesigningrequestsResource, opts))
}

// Create takes the representation of a certificateSigningRequest and creates it.  Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *FakeCertificateSigningRequests) Create(ctx context.Context, certificateSigningRequest *certificatesv1.CertificateSigningRequest, opts v1.CreateOptions) (result *certificatesv1.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(certificatesigningrequestsResource, certificateSigningRequest), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}

// Update takes the representation of a certificateSigningRequest and updates it. Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *FakeCertificateSigningRequests) Update(ctx context.Context, certificateSigningRequest *certificatesv1.CertificateSigningRequest, opts v1.UpdateOptions) (result *certificatesv1.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(certificatesigningrequestsResource, certificateSigningRequest), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCertificateSigningRequests) UpdateStatus(ctx context.Context, certificateSigningRequest *certificatesv1.CertificateSigningRequest, opts v1.UpdateOptions) (*certificatesv1.CertificateSigningRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(certificatesigningrequestsResource, "status", certificateSigningRequest), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}

// Delete takes name of the certificateSigningRequest and deletes it. Returns an error if one occurs.
func (c *FakeCertificateSigningRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(certificatesigningrequestsResource, name), &certificatesv1.CertificateSigningRequest{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCertificateSigningRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(certificatesigningrequestsResource, listOpts)

	_, err := c.Fake.Invokes(action, &certificatesv1.CertificateSigningRequestList{})
	return err
}

// Patch applies the patch and returns the patched certificateSigningRequest.
func (c *FakeCertificateSigningRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *certificatesv1.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(certificatesigningrequestsResource, name, pt, data, subresources...), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}

// UpdateApproval takes the representation of a certificateSigningRequest and updates it. Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *FakeCertificateSigningRequests) UpdateApproval(ctx context.Context, certificateSigningRequestName string, certificateSigningRequest *certificatesv1.CertificateSigningRequest, opts v1.UpdateOptions) (result *certificatesv1.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(certificatesigningrequestsResource, "approval", certificateSigningRequest), &certificatesv1.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*certificatesv1.CertificateSigningRequest), err
}
