package catalog

import (
	"strings"
	"testing"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
)

func TestRenderProductsMarkdownIsDeterministicAndMarksAmbiguousNames(t *testing.T) {
	products := []models.Product{
		{Name: "Premium VPS", Type: "vps", ServiceHandler: "vm6", Prices: map[string]float64{"month": 10, "hour": 1}},
		{Name: "Basic VPS", Type: "vps", ServiceHandler: "vm6", Prices: map[string]float64{"month": 5}},
		{Name: "premium vps", Type: "vps", ServiceHandler: "other", Prices: map[string]float64{"year": 100}},
		{Name: "Fast CPU", Type: "hicpu", ServiceHandler: "vm6", Prices: map[string]float64{"month": 20}},
		{Name: "Bare Metal", Type: "dedicated", ServiceHandler: "manual", Prices: map[string]float64{"month": 50}},
		{Name: "example", Type: "domain", ServiceHandler: "feru", Prices: map[string]float64{"year": 10}},
		{Name: "Archive", Type: "storage", ServiceHandler: "storage", Prices: map[string]float64{"month": 3}},
	}

	first := RenderProductsMarkdown(products)
	second := RenderProductsMarkdown([]models.Product{products[6], products[2], products[4], products[0], products[5], products[1], products[3]})
	if first != second {
		t.Fatal("catalog output depends on API result order")
	}
	if !strings.Contains(first, "`Basic VPS`") {
		t.Fatal("catalog does not contain product name")
	}
	if strings.Count(first, "⚠ ambiguous") != 2 {
		t.Fatalf("expected both duplicate names to be marked ambiguous:\n%s", first)
	}
	if !strings.Contains(first, "hour, month") {
		t.Fatal("payment terms are not sorted")
	}
	for _, expected := range []string{
		"| VPS | 3 |",
		"| HiCPU | 1 |",
		"| Dedicated | 1 |",
		"| Domains | 1 |",
		"| Other | 1 |",
		"## VPS (3)",
		"## HiCPU (1)",
		"## Dedicated (1)",
		"## Domains (1)",
		"## Other (1)",
	} {
		if !strings.Contains(first, expected) {
			t.Fatalf("catalog does not contain %q:\n%s", expected, first)
		}
	}
	if strings.Contains(first, "Generated at") {
		t.Fatal("catalog must not contain a volatile generation timestamp")
	}
}
