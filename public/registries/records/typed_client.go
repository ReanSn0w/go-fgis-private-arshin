package records

import (
	"context"
	"fmt"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

type Adapter[T Record] struct {
	Spec registries.Spec
}

type TypedClient[T Record] struct {
	arshin  *arshin.Client
	adapter Adapter[T]
}

type TypedPage[T Record] struct {
	TotalCount  int                     `json:"totalCount"`
	CurrentPage int                     `json:"currentPage"`
	PageSize    int                     `json:"pageSize"`
	Items       []T                     `json:"items"`
	RawPage     arshin.RegistryDataPage `json:"-"`
}

func NewTypedClient[T Record](client *arshin.Client, adapter Adapter[T]) *TypedClient[T] {
	return &TypedClient[T]{arshin: client, adapter: adapter}
}

func (c *TypedClient[T]) RegistrySpec() registries.Spec {
	return c.adapter.Spec
}

func (c *TypedClient[T]) List(ctx context.Context, query arshin.RegistryQuery) (*TypedPage[T], error) {
	if c.arshin == nil {
		return nil, fmt.Errorf("records: arshin client is nil")
	}

	raw, err := c.arshin.ListRegistryData(ctx, c.adapter.Spec.RegistryID, query)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(raw.Items))
	for _, item := range raw.Items {
		mapped, err := castRecord[T](MapRecord(c.adapter.Spec.RegistryID, item), c.adapter.Spec)
		if err != nil {
			return nil, err
		}
		items = append(items, mapped)
	}

	return &TypedPage[T]{
		TotalCount:  raw.TotalCount,
		CurrentPage: raw.CurrentPage,
		PageSize:    raw.PageSize,
		Items:       items,
		RawPage:     *raw,
	}, nil
}

func (c *TypedClient[T]) Get(ctx context.Context, itemID arshin.RegistryItemID) (T, error) {
	var zero T
	if c.arshin == nil {
		return zero, fmt.Errorf("records: arshin client is nil")
	}

	item, err := c.arshin.GetRegistryItem(ctx, c.adapter.Spec.RegistryID, itemID)
	if err != nil {
		return zero, err
	}
	return castRecord[T](MapItem(*item), c.adapter.Spec)
}

func castRecord[T Record](record Record, spec registries.Spec) (T, error) {
	typed, ok := record.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("records: registry %s mapped to %T, not requested type", spec.RegistryID, record)
	}
	return typed, nil
}

var (
	NormativeDocuments                 = Adapter[NormativeDocument]{Spec: registries.ND}
	CalibrationMarkCiphers             = Adapter[CalibrationMarkCipher]{Spec: registries.SCM}
	GSIStandards                       = Adapter[GSIStandard]{Spec: registries.GSI}
	InternationalTreaties              = Adapter[InternationalTreaty]{Spec: registries.MDG}
	InternationalComparisons           = Adapter[InternationalComparison]{Spec: registries.IC}
	ReferenceDataRecords               = Adapter[ReferenceData]{Spec: registries.SSD}
	StandardUnits                      = Adapter[StandardUnit]{Spec: registries.SU}
	PrimaryStandards                   = Adapter[PrimaryStandard]{Spec: registries.GPS}
	InternationalDocuments             = Adapter[InternationalDocument]{Spec: registries.MD}
	TechnicalInstrumentClassifications = Adapter[TechnicalInstrumentClassification]{Spec: registries.TSSI}
	RegulatedMeasurementLists          = Adapter[RegulatedMeasurementList]{Spec: registries.EPI}
	StandardSampleTypes                = Adapter[StandardSampleType]{Spec: registries.UTSO}
	ProductionNotices                  = Adapter[ProductionNotice]{Spec: registries.P1WF}
)
