// Copyright (c) 2024 Nikolai Osipov <nao99.dev@gmail.com>
//
// All rights are reserved
// Information about license can be found in the LICENSE file

package env

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Verifying correctness of env expansion", Label("env"), func() {
	When("an expansion is going", func() {
		AfterEach(func() {
			os.Clearenv()
		})

		Context("and env variable and default value both are presented in a content", func() {
			It("should successfully expand the env with its value when the env exists", func() {
				// given
				envVariableName := "CACHE_CONTROL"
				envDefaultValue := "no-cache"

				envValue := "no-store"

				_ = os.Setenv(envVariableName, envValue)

				content := fmt.Sprintf("${%s:%s}", envVariableName, envDefaultValue)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)
				expandedContentString := string(expandedContent)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContentString).To(Equal(envValue))
			})

			It("should successfully expand the env with the default value when the env doesn't exist", func() {
				// given
				envVariableName := "CACHE_CONTROL"
				envDefaultValue := "no-cache"

				content := fmt.Sprintf("${%s:%s}", envVariableName, envDefaultValue)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)
				expandedContentString := string(expandedContent)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContentString).To(Equal(envDefaultValue))
			})
		})

		Context("and only env variable is presented in a content", func() {
			It("should successfully expand the env with its value when the env exists", func() {
				// given
				envVariableName := "CACHE_CONTROL"
				envValue := "no-store"

				_ = os.Setenv(envVariableName, envValue)

				content := fmt.Sprintf("${%s}", envVariableName)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)
				expandedContentString := string(expandedContent)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContentString).To(Equal(envValue))
			})

			It("should return an error when the env doesn't exist", func() {
				// given
				envVariableName := "CACHE_CONTROL"

				content := fmt.Sprintf("${%s}", envVariableName)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)

				// then
				Expect(err).ToNot(BeNil())
				Expect(expandedContent).To(BeNil())
			})
		})

		Context("and env variable is presented with too long value", func() {
			It("should successfully expand the env with its value when the env exists #1", func() {
				// given
				envVariableName := "APP_DSN"
				envValue := "postgresql://localhost/app?user=dbuser&password=dbuser"

				_ = os.Setenv(envVariableName, envValue)

				content := fmt.Sprintf("${%s}", envVariableName)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)
				expandedContentString := string(expandedContent)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContentString).To(Equal(envValue))
			})

			It("should successfully expand the env with its value when the env exists #2", func() {
				// given
				dsnEnvVariableName := "APP_DSN"
				portEnvVariableName := "APP_PORT"

				dsnEnvValue := "postgresql://localhost/app?user=dbuser&password=dbuser"
				portEnvValue := "8080"

				_ = os.Setenv(dsnEnvVariableName, dsnEnvValue)
				_ = os.Setenv(portEnvVariableName, portEnvValue)

				contentString := `application:
				  name: "products-go"
				
				datasource:
				  dsn: ${%s}
				
				server:
				  port: ${%s}
				`

				content := fmt.Sprintf(contentString, dsnEnvVariableName, portEnvVariableName)
				contentBytes := []byte(content)

				// when
				_, err := ExpandEnvIn(contentBytes)

				// then
				Expect(err).To(BeNil())
			})
		})

		Context("and an env variable is presented in a content and the variable could be empty", func() {
			It("should successfully expand the env with an empty value the env doesn't exist", func() {
				// given
				envVariableName := "CACHE_CONTROL"

				content := fmt.Sprintf("${%s:}", envVariableName)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContent).To(BeEmpty())
			})
		})

		Context("and no env variable is presented", func() {
			It("should return a content as it is", func() {
				// given
				content := "simple content"
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)
				expandedContentString := string(expandedContent)

				// then
				Expect(err).To(BeNil())
				Expect(expandedContentString).To(Equal(content))
			})

			It("should return an error when the env doesn't exist", func() {
				// given
				envVariableName := "CACHE_CONTROL"

				content := fmt.Sprintf("${%s}", envVariableName)
				contentBytes := []byte(content)

				// when
				expandedContent, err := ExpandEnvIn(contentBytes)

				// then
				Expect(err).ToNot(BeNil())
				Expect(expandedContent).To(BeNil())
			})
		})

		Context("and content is nil", func() {
			It("should return an error", func() {
				// when
				expandedContent, err := ExpandEnvIn(nil)

				// then
				Expect(err).ToNot(BeNil())
				Expect(expandedContent).To(BeNil())
			})
		})
	})
})
