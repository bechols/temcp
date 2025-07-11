package main

import (
	"fmt"
	"os"

	"bechols/temcp/export"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Example usage: exporttool /path/to/export/file")
		os.Exit(1)
	}

	filename := os.Args[1]

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error reading file: %v\r\n", err)
		os.Exit(1)
	}

	workflows, err := export.DeserializeExportedWorkflows(bytes)
	if err != nil {
		fmt.Printf("error extracting workflow histories: %v\r\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully deserialized %d workflows \r\n", len(workflows.Items))

	for _, workflow := range workflows.Items {
		info, err := export.GetExportedWorkflowInformation(workflow)
		if err != nil {
			fmt.Printf("error extracting workflow information: %v\r\n", err)
			os.Exit(1)
		}

		fmt.Println(info)

		fmt.Println(export.FormatWorkflow(workflow))
		fmt.Println("----------------------------------------------------------")
		fmt.Println()
	}
}
