package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

var K8sClient *kubernetes.Clientset

func init() {
	config := &rest.Config{
		Host:        "110.41.142.160:8009",
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6Imt1SFFxZUNGMURxZXBXeGs0OFBSVkFWLUNnQXNEN0FDcU1URlJCQWY5ZncifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi1uenhsYiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjE5NDI5MjUxLWRkMDYtNGViMC05NDk0LWI1MDlmYjBkNTg5NyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.cVUTTbMYUObU7tGT8-bDa2HfCLYeLNyaHjVDGx_s0yWnm8qn7zfLkNeIGK6HVGRmevRikodLQBqq1aFUbiPiVwFFj1NkqEwHaAZKxJQabtjMexZkSIc-OoXiiIhi2f_ckr3mHPjlovqgjNXnHGeMQy-1XRNxCOe2-BpfeBHo_pKqAcvOtalmUDzkBQZ3qf1Qkuysfc37sLg1En0-00MzsWS-EUQ4Ea-k-wwwuc1jnjPbSeXQCb33pArTYcASPicWvJIjwgOSBobAFwBDbWvuqbX6vdc4rcegK3lQt7CTssLtjTZkDkpLR1yG6R9bMdL0EWOEJGZBemh44_FnMMHIzA",
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	K8sClient = c
}
