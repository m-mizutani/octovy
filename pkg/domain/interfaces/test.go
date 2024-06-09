package interfaces

/*
func FirestoreClientTest(t *testing.T, client Firestore) {
	type testData struct {
		Value string `firestore:"value"`
	}

	type testCase struct {
		get    []types.FireStoreRef
		put    []types.FireStoreRef
		input  testData
		output *testData
	}

	prefix := "test-" + uuid.NewString() + "-"
	p := func(id string) types.FSCollectionID {
		return types.FSCollectionID(prefix + id)
	}

	testCases := map[string]testCase{
		"simple put and get": {
			put: []types.FireStoreRef{{CollectionID: p("t1"), DocumentID: "doc1"}},
			get: []types.FireStoreRef{{CollectionID: p("t1"), DocumentID: "doc1"}},
			input: testData{
				Value: "hello",
			},
			output: &testData{
				Value: "hello",
			},
		},
		"nested put and get": {
			put: []types.FireStoreRef{{CollectionID: p("t2"), DocumentID: "doc1"}, {CollectionID: p("col2"), DocumentID: "doc2"}},
			get: []types.FireStoreRef{{CollectionID: p("t2"), DocumentID: "doc1"}, {CollectionID: p("col2"), DocumentID: "doc2"}},
			input: testData{
				Value: "hello",
			},
			output: &testData{
				Value: "hello",
			},
		},
		"nested put and get with different collection": {
			put: []types.FireStoreRef{{CollectionID: p("t3"), DocumentID: "doc1"}, {CollectionID: p("col2"), DocumentID: "doc3"}},
			get: []types.FireStoreRef{{CollectionID: p("t3"), DocumentID: "doc1"}, {CollectionID: p("col3"), DocumentID: "doc4"}},
			input: testData{
				Value: "hello",
			},
			output: nil,
		},
		"nested put and get with different document": {
			put: []types.FireStoreRef{{CollectionID: p("t4"), DocumentID: "doc1"}, {CollectionID: p("col2"), DocumentID: "doc5"}},
			get: []types.FireStoreRef{{CollectionID: p("t4"), DocumentID: "doc2"}, {CollectionID: p("col2"), DocumentID: "doc5"}},
			input: testData{
				Value: "hello",
			},
			output: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			if tc.put != nil {
				if err := client.Put(ctx, tc.input, tc.put...); err != nil {
					t.Fatal(err)
				}
			}

			if tc.get != nil {
				var output *testData
				if err := client.Get(ctx, &output, tc.get...); err != nil {
					t.Fatal(err)
				}

				if tc.output == nil && output != nil {
					t.Errorf("Unexpected output: %v", output)
				} else if tc.output != nil && output == nil {
					t.Errorf("No output")
				} else if tc.output != nil && output != nil {
					gt.Equal(t, tc.output.Value, output.Value)
				}
			}
		})
	}
}
*/
