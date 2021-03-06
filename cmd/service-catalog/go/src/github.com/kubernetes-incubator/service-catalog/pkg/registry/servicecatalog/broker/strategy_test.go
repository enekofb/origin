/*
Copyright 2016 The Kubernetes Authors.

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

package broker

import (
	"testing"

	sc "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func brokerWithOldSpec() *sc.ServiceBroker {
	return &sc.ServiceBroker{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Spec: sc.ServiceBrokerSpec{
			URL: "https://kubernetes.default.svc:443/brokers/template.k8s.io",
		},
		Status: sc.ServiceBrokerStatus{
			Conditions: []sc.ServiceBrokerCondition{
				{
					Type:   sc.ServiceBrokerConditionReady,
					Status: sc.ConditionFalse,
				},
			},
		},
	}
}

func brokerWithNewSpec() *sc.ServiceBroker {
	b := brokerWithOldSpec()
	b.Spec.URL = "new"
	return b
}

// TestServiceBrokerStrategyTrivial is the testing of the trivial hardcoded
// boolean flags.
func TestServiceBrokerStrategyTrivial(t *testing.T) {
	if brokerRESTStrategies.NamespaceScoped() {
		t.Errorf("broker create must not be namespace scoped")
	}
	if brokerRESTStrategies.NamespaceScoped() {
		t.Errorf("broker update must not be namespace scoped")
	}
	if brokerRESTStrategies.AllowCreateOnUpdate() {
		t.Errorf("broker should not allow create on update")
	}
	if brokerRESTStrategies.AllowUnconditionalUpdate() {
		t.Errorf("broker should not allow unconditional update")
	}
}

// TestBrokerCreate
func TestBroker(t *testing.T) {
	// Create a broker or brokers
	broker := &sc.ServiceBroker{
		Spec: sc.ServiceBrokerSpec{
			URL: "abcd",
		},
		Status: sc.ServiceBrokerStatus{
			Conditions: nil,
		},
	}

	// Canonicalize the broker
	brokerRESTStrategies.PrepareForCreate(nil, broker)

	if broker.Status.Conditions == nil {
		t.Fatalf("Fresh broker should have empty status")
	}
	if len(broker.Status.Conditions) != 0 {
		t.Fatalf("Fresh broker should have empty status")
	}
}

// TestBrokerUpdate tests that generation is incremented correctly when the
// spec of a Broker is updated.
func TestBrokerUpdate(t *testing.T) {
	cases := []struct {
		name                      string
		older                     *sc.ServiceBroker
		newer                     *sc.ServiceBroker
		shouldGenerationIncrement bool
	}{
		{
			name:  "no spec change",
			older: brokerWithOldSpec(),
			newer: brokerWithOldSpec(),
			shouldGenerationIncrement: false,
		},
		{
			name:  "spec change",
			older: brokerWithOldSpec(),
			newer: brokerWithNewSpec(),
			shouldGenerationIncrement: true,
		},
	}

	for i := range cases {
		brokerRESTStrategies.PrepareForUpdate(nil, cases[i].newer, cases[i].older)

		if cases[i].shouldGenerationIncrement {
			if e, a := cases[i].older.Generation+1, cases[i].newer.Generation; e != a {
				t.Fatalf("%v: expected %v, got %v for generation", cases[i].name, e, a)
			}
		} else {
			if e, a := cases[i].older.Generation, cases[i].newer.Generation; e != a {
				t.Fatalf("%v: expected %v, got %v for generation", cases[i].name, e, a)
			}
		}
	}
}
