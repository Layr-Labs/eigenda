package relay

// var (
// 	logger              = test.GetLogger()
// 	localstackContainer *testbed.LocalStackContainer
// 	UUID                = uuid.New()
// 	metadataTableName   = fmt.Sprintf("test-BlobMetadata-%v", UUID)
// 	prover              *p.Prover
// 	bucketName          = fmt.Sprintf("test-bucket-%v", UUID)
// )

// const (
// 	localstackPort = "4570"
// 	localstackHost = "http://0.0.0.0:4570"
// )

// func setup(t *testing.T) {
// 	ctx := t.Context()
// 	deployLocalStack := (os.Getenv("DEPLOY_LOCALSTACK") != "false")

// 	_, b, _, _ := runtime.Caller(0)
// 	rootPath := filepath.Join(filepath.Dir(b), "..")
// 	changeDirectory(filepath.Join(rootPath, "inabox"))

// 	if deployLocalStack {
// 		var err error
// 		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
// 			ExposeHostPort: true,
// 			HostPort:       localstackPort,
// 			Services:       []string{"s3", "dynamodb"},
// 			Logger:         logger,
// 		})
// 		require.NoError(t, err)
// 	}

// 	// Only set up the prover once, it's expensive
// 	if prover == nil {
// 		config := &kzg.KzgConfig{
// 			G1Path:          "../resources/srs/g1.point",
// 			G2Path:          "../resources/srs/g2.point",
// 			CacheDir:        "../resources/srs/SRSTables",
// 			SRSOrder:        8192,
// 			SRSNumberToLoad: 8192,
// 			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
// 			LoadG2Points:    true,
// 		}
// 		var err error
// 		prover, err = p.NewProver(config, nil)
// 		require.NoError(t, err)
// 	}
// }

// func changeDirectory(path string) {
// 	err := os.Chdir(path)
// 	if err != nil {
// 		logger.Fatal("Failed to change directories. Error: ", err)
// 	}

// 	newDir, err := os.Getwd()
// 	if err != nil {
// 		logger.Fatal("Failed to get working directory. Error: ", err)
// 	}
// 	logger.Debug("Current Working Directory: %s", newDir)
// }

// func teardown(t *testing.T) {
// 	t.Helper()
// 	deployLocalStack := (os.Getenv("DEPLOY_LOCALSTACK") != "false")

// 	if deployLocalStack {
// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()
// 		_ = localstackContainer.Terminate(ctx)
// 	}
// }

// func buildMetadataStore(t *testing.T) *blobstore.BlobMetadataStore {
// 	t.Helper()
// 	ctx := t.Context()

// 	err := os.Setenv("AWS_ACCESS_KEY_ID", "localstack")
// 	require.NoError(t, err)
// 	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "localstack")
// 	require.NoError(t, err)

// 	cfg := aws.ClientConfig{
// 		Region:          "us-east-1",
// 		AccessKey:       "localstack",
// 		SecretAccessKey: "localstack",
// 		EndpointURL:     localstackHost,
// 	}

// 	_, err = test_utils.CreateTable(
// 		ctx,
// 		cfg,
// 		metadataTableName,
// 		blobstore.GenerateTableSchema(metadataTableName, 10, 10))
// 	if err != nil {
// 		if !strings.Contains(err.Error(), "ResourceInUseException: Table already exists") {
// 			require.NoError(t, err)
// 		}
// 	}

// 	dynamoClient, err := dynamodb.NewClient(cfg, logger)
// 	require.NoError(t, err)

// 	return blobstore.NewBlobMetadataStore(
// 		dynamoClient,
// 		logger,
// 		metadataTableName)
// }

// func buildBlobStore(t *testing.T, logger logging.Logger) *blobstore.BlobStore {
// 	t.Helper()
// 	ctx := t.Context()

// 	cfg := aws.DefaultClientConfig()
// 	cfg.Region = "us-east-1"
// 	cfg.AccessKey = "localstack"
// 	cfg.SecretAccessKey = "localstack"
// 	cfg.EndpointURL = localstackHost

// 	client, err := s3.NewClient(ctx, *cfg, logger)
// 	require.NoError(t, err)

// 	err = client.CreateBucket(ctx, bucketName)
// 	require.NoError(t, err)

// 	return blobstore.NewBlobStore(bucketName, client, logger)
// }

// func buildChunkStore(t *testing.T, logger logging.Logger) (chunkstore.ChunkReader, chunkstore.ChunkWriter) {
// 	t.Helper()
// 	ctx := t.Context()

// 	cfg := aws.ClientConfig{
// 		Region:          "us-east-1",
// 		AccessKey:       "localstack",
// 		SecretAccessKey: "localstack",
// 		EndpointURL:     localstackHost,
// 	}

// 	client, err := s3.NewClient(ctx, cfg, logger)
// 	require.NoError(t, err)

// 	err = client.CreateBucket(ctx, bucketName)
// 	require.NoError(t, err)

// 	// intentionally use very small fragment size
// 	chunkWriter := chunkstore.NewChunkWriter(logger, client, bucketName, 32)
// 	chunkReader := chunkstore.NewChunkReader(logger, client, bucketName)

// 	return chunkReader, chunkWriter
// }

// func newMockChainReader(t *testing.T) *coremock.MockWriter {
// 	t.Helper()
// 	w := &coremock.MockWriter{}
// 	w.On("GetAllVersionedBlobParams", mock.Anything).Return(mockBlobParamsMap(t), nil)
// 	return w
// }

// func mockBlobParamsMap(t *testing.T) map[v2.BlobVersion]*core.BlobVersionParameters {
// 	t.Helper()
// 	blobParams := &core.BlobVersionParameters{
// 		NumChunks:       8192,
// 		CodingRate:      8,
// 		MaxNumOperators: 2048,
// 	}

// 	return map[v2.BlobVersion]*core.BlobVersionParameters{
// 		0: blobParams,
// 	}
// }

// func randomBlob(t *testing.T) (*v2.BlobHeader, []byte) {
// 	t.Helper()

// 	data := random.RandomBytes(225)

// 	data = codec.ConvertByPaddingEmptyByte(data)
// 	commitments, err := prover.GetCommitmentsForPaddedLength(data)
// 	require.NoError(t, err)
// 	require.NoError(t, err)
// 	commitmentProto, err := commitments.ToProtobuf()
// 	require.NoError(t, err)

// 	blobHeaderProto := &pbcommonv2.BlobHeader{
// 		Version:       0,
// 		QuorumNumbers: []uint32{0, 1},
// 		Commitment:    commitmentProto,
// 		PaymentHeader: &pbcommonv2.PaymentHeader{
// 			AccountId:         gethcommon.BytesToAddress(random.RandomBytes(20)).Hex(),
// 			Timestamp:         5,
// 			CumulativePayment: big.NewInt(100).Bytes(),
// 		},
// 	}
// 	blobHeader, err := v2.BlobHeaderFromProtobuf(blobHeaderProto)
// 	require.NoError(t, err)

// 	return blobHeader, data
// }

// func randomBlobChunks(t *testing.T) (*v2.BlobHeader, []byte, []*encoding.Frame) {
// 	t.Helper()
// 	header, data := randomBlob(t)

// 	params := encoding.ParamsFromMins(16, 16)
// 	_, frames, err := prover.EncodeAndProve(data, params)
// 	require.NoError(t, err)

// 	return header, data, frames
// }

// func disassembleFrames(t *testing.T, frames []*encoding.Frame) ([]rs.FrameCoeffs, []*encoding.Proof) {
// 	t.Helper()
// 	rsFrames := make([]rs.FrameCoeffs, len(frames))
// 	proofs := make([]*encoding.Proof, len(frames))

// 	for i, frame := range frames {
// 		rsFrames[i] = frame.Coeffs
// 		proofs[i] = &frame.Proof
// 	}

// 	return rsFrames, proofs
// }
