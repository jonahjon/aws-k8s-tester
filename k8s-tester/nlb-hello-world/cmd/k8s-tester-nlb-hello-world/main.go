// k8s-tester-nlb-hello-world installs Kubernetes NLB hello world tester.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-k8s-tester/client"
	nlb_hello_world "github.com/aws/aws-k8s-tester/k8s-tester/nlb-hello-world"
	"github.com/aws/aws-k8s-tester/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:        "k8s-tester-nlb-hello-world",
	Short:      "Kubernetes NLB hello world tester",
	SuggestFor: []string{"nlb-hello-world"},
}

func init() {
	cobra.EnablePrefixMatching = true
}

var (
	enablePrompt   bool
	logLevel       string
	logOutputs     []string
	namespace      string
	kubectlPath    string
	kubeConfigPath string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&enablePrompt, "enable-prompt", true, "'true' to enable prompt mode")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", log.DefaultLogLevel, "Logging level")
	rootCmd.PersistentFlags().StringSliceVar(&logOutputs, "log-outputs", []string{"stderr"}, "Additional logger outputs")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "test-namespace", "'true' to auto-generate path for create config/cluster, overwrites existing --path value")
	rootCmd.PersistentFlags().StringVar(&kubectlPath, "kubectl-path", "", "kubectl path")
	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeconfig-path", "", "KUBECONFIG path")

	rootCmd.AddCommand(
		newApply(),
		newDelete(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "k8s-tester-nlb-hello-world failed %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

var (
	deploymentNodeSelector string
	deploymentReplicas     int32
)

func newApply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply tests",
		Run:   createApplyFunc,
	}
	cmd.PersistentFlags().StringVar(&deploymentNodeSelector, "deployment-node-selector", "", "map of deployment node selector, must be valid JSON format")
	cmd.PersistentFlags().Int32Var(&deploymentReplicas, "deployment-replicas", 2, "number of deployment replicas")
	return cmd
}

func createApplyFunc(cmd *cobra.Command, args []string) {
	lg, logWriter, _, err := log.NewWithStderrWriter(logLevel, logOutputs)
	if err != nil {
		panic(err)
	}
	_ = zap.ReplaceGlobals(lg)

	nodeSelector := make(map[string]string)
	if err := json.Unmarshal([]byte(deploymentNodeSelector), &nodeSelector); err != nil {
		panic(fmt.Errorf("failed to parse %q (%v)", deploymentNodeSelector, err))
	}

	cfg := nlb_hello_world.Config{
		EnablePrompt: enablePrompt,
		Logger:       lg,
		LogWriter:    logWriter,
		Namespace:    namespace,
		ClientConfig: &client.Config{
			Logger:         lg,
			KubectlPath:    kubectlPath,
			KubeConfigPath: kubeConfigPath,
		},
		DeploymentNodeSelector: nodeSelector,
		DeploymentReplicas:     deploymentReplicas,
	}

	ts := nlb_hello_world.New(cfg)
	if err := ts.Apply(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to apply (%v)\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n*********************************\n")
	fmt.Printf("'k8s-tester-nlb-hello-world apply' success\n")
}

func newDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Run:   createDeleteFunc,
	}
	return cmd
}

func createDeleteFunc(cmd *cobra.Command, args []string) {
	lg, logWriter, _, err := log.NewWithStderrWriter(logLevel, logOutputs)
	if err != nil {
		panic(err)
	}
	_ = zap.ReplaceGlobals(lg)

	cfg := nlb_hello_world.Config{
		EnablePrompt: enablePrompt,
		Logger:       lg,
		LogWriter:    logWriter,
		Namespace:    namespace,
		ClientConfig: &client.Config{
			Logger:         lg,
			KubectlPath:    kubectlPath,
			KubeConfigPath: kubeConfigPath,
		},
	}

	ts := nlb_hello_world.New(cfg)
	if err := ts.Delete(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete (%v)\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n*********************************\n")
	fmt.Printf("'k8s-tester-nlb-hello-world delete' success\n")
}
