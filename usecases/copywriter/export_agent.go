package copywriter

import (
	"fmt"
	"os"

	"github.com/nullswan/nomi/internal/tools"
)

// This agent is responsible for exporting the final content
// Capabilities include: Export to file, export to different formats
type exportAgent struct {
	logger  tools.Logger
	project string
}

// TODO(nullswan): add project handling
// TODO(nullswan): add file manager tools
func newExportAgent(
	logger tools.Logger,
	project string,
) *exportAgent {
	return &exportAgent{
		logger:  logger,
		project: project,
	}
}

func (e *exportAgent) ExportToFile(prefix, content string) error {
	// TODO(nullswan): Encode prefix
	fileName := e.project + "-" + prefix + ".txt"

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Printf("Content exported to %s\n", fileName)
	return nil
}
