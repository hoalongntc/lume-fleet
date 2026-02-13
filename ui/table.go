package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/hoalong/lume-fleet/fleet"
	"github.com/hoalong/lume-fleet/lume"
)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	bold   = lipgloss.NewStyle().Bold(true)
	dimBorder = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

// StatusRow holds display data for one VM in the status table.
type StatusRow struct {
	Name   string
	State  string
	IP     string
	OS     string
	CPU    int
	Memory string
	Tags   []string
}

// BuildStatusRows merges resolved fleet VMs with actual Lume state.
func BuildStatusRows(resolved []fleet.ResolvedVM, actual []lume.VM) []StatusRow {
	index := make(map[string]lume.VM, len(actual))
	for _, vm := range actual {
		index[vm.Name] = vm
	}

	var rows []StatusRow
	for _, r := range resolved {
		row := StatusRow{
			Name:   r.Name,
			OS:     r.OS,
			CPU:    r.CPU,
			Memory: r.Memory,
			Tags:   r.Tags,
		}

		if vm, ok := index[r.Name]; ok {
			row.State = vm.Status
			if vm.IPAddress != nil {
				row.IP = *vm.IPAddress
			}
			row.CPU = vm.CPUCount
			row.Memory = formatBytes(vm.MemorySize)
		} else {
			row.State = "not created"
		}

		if row.IP == "" {
			row.IP = "-"
		}

		rows = append(rows, row)
	}

	return rows
}

// RenderStatusTable outputs a formatted status table.
func RenderStatusTable(rows []StatusRow, macosRunning int) string {
	var sb strings.Builder

	header := fmt.Sprintf("  Fleet Status (%d VMs)  |  macOS: %d/2 slots", len(rows), macosRunning)
	sb.WriteString(bold.Render(header))
	sb.WriteString("\n\n")

	tableRows := make([][]string, len(rows))
	for i, r := range rows {
		tableRows[i] = []string{
			r.Name,
			colorizeState(r.State),
			r.IP,
			r.OS,
			fmt.Sprintf("%d", r.CPU),
			r.Memory,
			strings.Join(r.Tags, ", "),
		}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(dimBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		}).
		Headers("NAME", "STATE", "IP", "OS", "CPU", "MEMORY", "TAGS").
		Rows(tableRows...)

	sb.WriteString(t.String())
	return sb.String()
}

func colorizeState(s string) string {
	switch strings.ToLower(s) {
	case "running":
		return green.Render("running")
	case "stopped":
		return yellow.Render("stopped")
	default:
		return gray.Render(s)
	}
}

func formatBytes(b int64) string {
	gb := float64(b) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.0fGB", gb)
	}
	mb := float64(b) / (1024 * 1024)
	return fmt.Sprintf("%.0fMB", mb)
}
