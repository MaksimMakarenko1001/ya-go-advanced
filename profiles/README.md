# Сравнение профилей

```
go tool pprof -sample_index=alloc_objects -top -diff_base=profiles/base.pprof profiles/result.pprof

File: main
Type: alloc_objects
Time: 2026-03-15 22:46:54 MSK
Showing nodes accounting for -47399714, 64.63% of 73342163 total
Dropped 402 nodes (cum <= 366710)
      flat  flat%   sum%        cum   cum%
  -1671244  2.28%  2.28%   -1671244  2.28%  time.Time.MarshalJSON
  -1473463  2.01%  4.29%   -4897005  6.68%  encoding/json.Marshal
  -1306486  1.78%  6.07%   -1306486  1.78%  reflect.growslice
  -1247073  1.70%  7.77%   -1247073  1.70%  sync.(*Pool).pinSlow
  -1245511  1.70%  9.47%   -1409353  1.92%  net/textproto.readMIMEHeader
  -1201580  1.64% 11.11%   -2029075  2.77%  crypto/internal/fips140/hmac.New[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
  -1200280  1.64% 12.74%   -1200280  1.64%  compress/flate.newHuffmanEncoder (inline)
  -1114128  1.52% 14.26%   -1572887  2.14%  encoding/json.(*decodeState).literalStore
  -1084616  1.48% 15.74%   -1084616  1.48%  bytes.growSlice
   -917554  1.25% 16.99%   -2442573  3.33%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/logger.(*ZapLogger).LogHTTP
   -915392  1.25% 18.24%  -29121729 39.71%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateBatchService/v0.(*Service).Do
   -912058  1.24% 19.48%   -1742196  2.38%  github.com/jackc/pgx/v5/pgproto3.(*ParameterStatus).Decode
   -901133  1.23% 20.71%    -901133  1.23%  encoding/json.(*scanner).pushParseState
   -830138  1.13% 21.84%    -830138  1.13%  bytes.(*Buffer).ReadBytes (inline)
   -828822  1.13% 22.97%   -2029102  2.77%  compress/flate.newHuffmanBitWriter (inline)
   -827495  1.13% 24.10%    -827495  1.13%  crypto/internal/fips140/sha256.New (inline)
   -779803  1.06% 25.17%    -999840  1.36%  context.(*cancelCtx).propagateCancel
   -754799  1.03% 26.19%    -754799  1.03%  vendor/golang.org/x/crypto/cryptobyte.(*String).ReadASN1ObjectIdentifier
   -737303  1.01% 27.20%    -737303  1.01%  strings.(*Builder).grow
   -734060  1.00% 28.20%   -1733900  2.36%  context.withCancel (inline)
   -699067  0.95% 29.15%  -26505444 36.14%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/storage/pg.(*Repository).AddUpdateBatch
   -682687  0.93% 30.08%    -988539  1.35%  crypto/x509.parseName
   -672966  0.92% 31.00%  -48850232 66.61%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.WithHash-fm.API.WithHash.func1
   -661393   0.9% 31.90%    -759698  1.04%  net/textproto.MIMEHeader.Set (inline)
   -653202  0.89% 32.79%    -653202  0.89%  bufio.NewReaderSize (inline)
   -644463  0.88% 33.67%    -644463  0.88%  os.newFile
   -644452  0.88% 34.55%   -1051935  1.43%  github.com/jackc/pgx/v5/stdlib.(*Rows).Next
   -632090  0.86% 35.41%    -833338  1.14%  compress/flate.NewReader
   -595295  0.81% 36.23%    -595295  0.81%  syscall.ByteSliceFromString
   -589833   0.8% 37.03%   -2323733  3.17%  context.WithCancel
   -589832   0.8% 37.83%    -589832   0.8%  net/textproto.canonicalMIMEHeaderKey
   -577438  0.79% 38.62%    -577438  0.79%  compress/flate.(*huffmanEncoder).generate
   -571160  0.78% 39.40%    -571160  0.78%  context.(*cancelCtx).Done
   -557064  0.76% 40.16%    -557064  0.76%  github.com/jackc/pgx/v5/pgconn.(*PgConn).convertRowDescription
   -524296  0.71% 40.87%   -1074371  1.46%  database/sql.(*Rows).close
   -507480  0.69% 41.57%    -507480  0.69%  net/http.Header.Clone (inline)
   -494945  0.67% 42.24%   -3015678  4.11%  crypto/x509.parseCertificate
   -480620  0.66% 42.90%    -480620  0.66%  context.WithValue
   -480607  0.66% 43.55%    -480607  0.66%  io.LimitReader (inline)
   -476283  0.65% 44.20%   -5195862  7.08%  net/http.(*conn).readRequest
   -465055  0.63% 44.84%    -465055  0.63%  bufio.NewWriterSize (inline)
   -458766  0.63% 45.46%   -1836807  2.50%  encoding/json.(*decodeState).object
   -458759  0.63% 46.09%    -458759  0.63%  reflect.New
   -445679  0.61% 46.69%   -5517631  7.52%  github.com/jackc/pgx/v5/stdlib.(*Conn).QueryContext
   -442381   0.6% 47.30%   -3958153  5.40%  crypto/x509.(*CertPool).AppendCertsFromPEM
   -434203  0.59% 47.89%    -434203  0.59%  syscall.anyToSockaddr
   -434192  0.59% 48.48%  -57094086 77.85%  net/http.(*conn).serve
   -409620  0.56% 49.04%  -43836503 59.77%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.WithLogging-fm.API.WithLogging.func1
   -407742  0.56% 49.60%   -2673141  3.64%  github.com/jackc/pgx/v5/pgtype.(*encodePlanJSONCodecEitherFormatMarshal).Encode
   -403111  0.55% 50.15%    -403111  0.55%  compress/flate.(*compressor).initDeflate (inline)
   -393225  0.54% 50.68%    -393225  0.54%  go.uber.org/zap.ByteString (inline)
   -371388  0.51% 51.19%    -371388  0.51%  encoding/base64.(*Encoding).EncodeToString
   -371379  0.51% 51.69%  -39614463 54.01%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.RegisterHandlers.func4.DoUpdateBatchJSONResponse.1
   -363196   0.5% 52.19%    -455206  0.62%  encoding/json.newEncodeState
   -344074  0.47% 52.66%    -344074  0.47%  crypto/internal/fips140/sha256.(*Digest).Sum
   -327689  0.45% 53.11%    -365142   0.5%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Watch
   -327685  0.45% 53.55%    -327685  0.45%  net/http.(*connReader).startBackgroundRead
   -327684  0.45% 54.00%    -327684  0.45%  net.JoinHostPort (inline)
   -319527  0.44% 54.43%    -319527  0.44%  net.newFD (inline)
   -314477  0.43% 54.86%    -314477  0.43%  io.init.func1
   -311315  0.42% 55.29%   -1395931  1.90%  bytes.(*Buffer).grow
   -307947  0.42% 55.71%   -6628757  9.04%  github.com/jackc/pgx/v5/pgconn.connectOne
   -305841  0.42% 56.12%  -47849581 65.24%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.MiddlewareCompress.func1
   -298254  0.41% 56.53%    -298254  0.41%  github.com/jackc/pgx/v5/pgtype.scanPlanAnyTextToBytes.Scan
   -294913   0.4% 56.93%   -6923670  9.44%  github.com/jackc/pgx/v5/pgconn.connectPreferred
   -286738  0.39% 57.32%    -592579  0.81%  net/http.readTransfer
   -262148  0.36% 57.68%  -21548138 29.38%  database/sql.(*DB).query
   -245797  0.34% 58.02%   -7004860  9.55%  database/sql.(*DB).queryDC
   -240304  0.33% 58.34%    -404154  0.55%  net/textproto.(*Reader).ReadLine (inline)
   -238650  0.33% 58.67%    -391570  0.53%  os.statNolog
   -217974   0.3% 58.97%   -3193167  4.35%  net/http.readRequest
   -213005  0.29% 59.26%   -1021853  1.39%  database/sql.(*Rows).initContextClose
   -201325  0.27% 59.53%   -2633568  3.59%  compress/flate.NewWriter (inline)
   -196613  0.27% 59.80%  -26094933 35.58%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db.(*PGConnect).QueryWithOneResultJSON
   -195305  0.27% 60.07%   -2070662  2.82%  github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage
   -193759  0.26% 60.33%   -1766712  2.41%  compress/gzip.NewReader (inline)
   -188424  0.26% 60.59%    -624562  0.85%  github.com/jackc/pgx/v5/pgconn.buildConnectOneConfigs
   -185694  0.25% 60.84%    -305848  0.42%  github.com/jackc/pgx/v5/internal/stmtcache.NewLRUCache (inline)
   -181640  0.25% 61.09%   -5747866  7.84%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions
   -180230  0.25% 61.33%   -1545212  2.11%  net/http.(*Server).Serve
   -169283  0.23% 61.56%   -8217250 11.20%  github.com/jackc/pgx/v5.connect
   -156579  0.21% 61.78%   -1376688  1.88%  encoding/json.Unmarshal
   -152921  0.21% 61.99%    -452317  0.62%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/storage/inmemory.(*Repository).Save
   -141998  0.19% 62.18%    -710022  0.97%  github.com/jackc/pgx/v5/pgconn.computeClientProof
   -131074  0.18% 62.36%    -364600   0.5%  fmt.Sprint
   -131074  0.18% 62.54%   -1776424  2.42%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/hashService/v0.(*Service).Hash
   -124521  0.17% 62.71%    -353911  0.48%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Prepare
   -117036  0.16% 62.87%   -2405909  3.28%  github.com/jackc/pgx/v5/pgconn.(*PgConn).scramAuth
   -109229  0.15% 63.02%   -1434260  1.96%  github.com/jackc/pgx/v5/pgconn.(*scramClient).clientFinalMessage
   -109229  0.15% 63.17%    -407483  0.56%  github.com/jackc/pgx/v5/stdlib.(*Rows).Next.func10
    -98307  0.13% 63.30%    -325561  0.44%  crypto/internal/fips140/pbkdf2.Key[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
    -98306  0.13% 63.43%   -1379620  1.88%  database/sql.(*Rows).nextLocked
    -96133  0.13% 63.56%   -4309928  5.88%  github.com/jackc/pgx/v5/pgconn.configTLS
    -94674  0.13% 63.69%    -728736  0.99%  github.com/jackc/pgx/v5/pgconn.defaultSettings
    -87398  0.12% 63.81%    -912459  1.24%  net.(*Dialer).DialContext
    -86029  0.12% 63.93%  -14180190 19.33%  github.com/jackc/pgx/v5/stdlib.(*driverConnector).Connect
    -76460   0.1% 64.03%    -450594  0.61%  github.com/jackc/pgx/v5/pgconn.computeServerSignature
    -71003 0.097% 64.13%   -3155623  4.30%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).appendParam
    -65539 0.089% 64.22%    -426158  0.58%  fmt.Errorf
    -65537 0.089% 64.31%    -466122  0.64%  fmt.(*pp).doPrintf
    -52433 0.071% 64.38%    -562907  0.77%  net.(*sysDialer).dialSingle
    -45062 0.061% 64.44%  -14292053 19.49%  database/sql.(*DB).conn
    -37817 0.052% 64.49%    -391728  0.53%  github.com/jackc/pgx/v5.(*Conn).Prepare
    -32772 0.045% 64.54%   -3084620  4.21%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).encodeExtendedParamValue
    -21846  0.03% 64.57%    -480479  0.66%  github.com/jackc/pgx/v5/pgconn.connectOne.func1 (inline)
    -14750  0.02% 64.59%   -5762616  7.86%  github.com/jackc/pgx/v5.ParseConfigWithOptions
    -10923 0.015% 64.60%    -299479  0.41%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/outbox.(*Repository).OutboxGetNext
    -10923 0.015% 64.62%   -1032224  1.41%  os.openFileNolog
     -4681 0.0064% 64.62%   -1520428  2.07%  runtime.main
     -2344 0.0032% 64.63%    -303532  0.41%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/auditRemoteService/v0.(*Service).Do
       -14 1.9e-05% 64.63%   -2432243  3.32%  compress/flate.(*compressor).init
         0     0% 64.63%    -327685  0.45%  bufio.(*Reader).Read
         0     0% 64.63%    -875840  1.19%  bufio.(*Writer).Flush
         0     0% 64.63%    -653202  0.89%  bufio.NewReader (inline)
         0     0% 64.63%    -836619  1.14%  bytes.(*Buffer).Write
         0     0% 64.63%    -802995  1.09%  compress/flate.(*Writer).Close (inline)
         0     0% 64.63%    -802995  1.09%  compress/flate.(*compressor).close
         0     0% 64.63%    -577451  0.79%  compress/flate.(*compressor).deflate
         0     0% 64.63%    -577451  0.79%  compress/flate.(*compressor).writeBlock
         0     0% 64.63%    -386555  0.53%  compress/flate.(*huffmanBitWriter).indexTokens
         0     0% 64.63%    -577451  0.79%  compress/flate.(*huffmanBitWriter).writeBlock
         0     0% 64.63%   -1572953  2.14%  compress/gzip.(*Reader).Reset
         0     0% 64.63%   -1161023  1.58%  compress/gzip.(*Reader).readHeader
         0     0% 64.63%    -817334  1.11%  compress/gzip.(*Writer).Close
         0     0% 64.63%   -2789241  3.80%  compress/gzip.(*Writer).Write
         0     0% 64.63%   -1872046  2.55%  crypto/hmac.New
         0     0% 64.63%    -749661  1.02%  crypto/hmac.New.UnwrapNew[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }].func1
         0     0% 64.63%    -344074  0.47%  crypto/internal/fips140/hmac.(*HMAC).Sum
         0     0% 64.63%    -325561  0.44%  crypto/pbkdf2.Key[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
         0     0% 64.63%    -827495  1.13%  crypto/sha256.New
         0     0% 64.63%   -3015678  4.11%  crypto/x509.ParseCertificate
         0     0% 64.63%    -370295   0.5%  crypto/x509.parseExtension
         0     0% 64.63%  -21548138 29.38%  database/sql.(*DB).QueryContext
         0     0% 64.63%  -21548138 29.38%  database/sql.(*DB).QueryContext.func1
         0     0% 64.63%  -21548138 29.38%  database/sql.(*DB).QueryRowContext (inline)
         0     0% 64.63%    -391986  0.53%  database/sql.(*DB).putConn
         0     0% 64.63%   -5737210  7.82%  database/sql.(*DB).queryDC.func1
         0     0% 64.63%  -21686734 29.57%  database/sql.(*DB).retry
         0     0% 64.63%   -2659531  3.63%  database/sql.(*Row).Scan
         0     0% 64.63%   -1014295  1.38%  database/sql.(*Rows).Close
         0     0% 64.63%   -1379620  1.88%  database/sql.(*Rows).Next
         0     0% 64.63%   -1379620  1.88%  database/sql.(*Rows).Next.func1
         0     0% 64.63%    -391986  0.53%  database/sql.(*driverConn).Close
         0     0% 64.63%    -391986  0.53%  database/sql.(*driverConn).releaseConn
         0     0% 64.63%   -5517631  7.52%  database/sql.ctxDriverQuery
         0     0% 64.63%   -7690810 10.49%  database/sql.withLock
         0     0% 64.63%    -540689  0.74%  encoding/asn1.ObjectIdentifier.String
         0     0% 64.63%   -3311839  4.52%  encoding/json.(*Decoder).Decode
         0     0% 64.63%    -716902  0.98%  encoding/json.(*Decoder).readValue
         0     0% 64.63%   -3733126  5.09%  encoding/json.(*decodeState).array
         0     0% 64.63%    -327685  0.45%  encoding/json.(*decodeState).scanWhile
         0     0% 64.63%   -3765894  5.13%  encoding/json.(*decodeState).unmarshal
         0     0% 64.63%   -3733126  5.09%  encoding/json.(*decodeState).value
         0     0% 64.63%   -2953698  4.03%  encoding/json.(*encodeState).marshal
         0     0% 64.63%   -2953698  4.03%  encoding/json.(*encodeState).reflectValue
         0     0% 64.63%   -1530746  2.09%  encoding/json.addrMarshalerEncoder
         0     0% 64.63%    -321692  0.44%  encoding/json.appendCompact
         0     0% 64.63%   -2230805  3.04%  encoding/json.arrayEncoder.encode
         0     0% 64.63%   -2120507  2.89%  encoding/json.condAddrEncoder.encode
         0     0% 64.63%    -458759  0.63%  encoding/json.indirect
         0     0% 64.63%    -589761   0.8%  encoding/json.marshalerEncoder
         0     0% 64.63%    -321985  0.44%  encoding/json.newScanner
         0     0% 64.63%   -2230805  3.04%  encoding/json.sliceEncoder.encode
         0     0% 64.63%    -901133  1.23%  encoding/json.stateBeginValue
         0     0% 64.63%    -491528  0.67%  encoding/json.stateBeginValueOrEmpty
         0     0% 64.63%    -374137  0.51%  encoding/json.stringEncoder
         0     0% 64.63%   -2863572  3.90%  encoding/json.structEncoder.encode
         0     0% 64.63%    -466122  0.64%  fmt.(*pp).printArg
         0     0% 64.63%   -1540531  2.10%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/config.(*DI).Start
         0     0% 64.63%  -24521632 33.43%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db.(*PGConnect).QueryWithOneResult
         0     0% 64.63%  -24521632 33.43%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db.(*PGConnect).QueryWithOneResult.(*PGConnect).QueryWithOneResult.(*Backoff).WithLinear.func2.func3 (inline)
         0     0% 64.63%  -24207669 33.01%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db.(*PGConnect).QueryWithOneResult.func1 (inline)
         0     0% 64.63%    -817326  1.11%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*compressWriter).Close
         0     0% 64.63%   -2789191  3.80%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*compressWriter).Write
         0     0% 64.63%   -2644357  3.61%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*compressWriter).WriteHeader
         0     0% 64.63%    -395561  0.54%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*responseHashWriter).Write
         0     0% 64.63%   -2447746  3.34%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*responseHashWriter).WriteHeader
         0     0% 64.63%   -2992629  4.08%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*responseWriter).Write
         0     0% 64.63%   -2644357  3.61%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.(*responseWriter).WriteHeader
         0     0% 64.63%  -44749955 61.02%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.ServeHTTP
         0     0% 64.63%   -1776424  2.42%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.WithHash-fm.API.WithHash.func1.1
         0     0% 64.63%  -40984310 55.88%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.API.WithSync-fm.API.WithSync.func1
         0     0% 64.63%   -6029943  8.22%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler.WriteJSONResult
         0     0% 64.63%    -299396  0.41%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/encode.(*JSONEncode).Encode
         0     0% 64.63%   -1369847  1.87%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/dumpMetricService/v0.(*Service).WriteDump
         0     0% 64.63%    -430961  0.59%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/worker/sworker.(*SimpleWorker).doJob
         0     0% 64.63%    -430961  0.59%  github.com/MaksimMakarenko1001/ya-go-advanced/internal/worker/sworker.(*SimpleWorker).run.func1
         0     0% 64.63%   -1488673  2.03%  github.com/MaksimMakarenko1001/ya-go-advanced/pkg.JSONMust (inline)
         0     0% 64.63%  -43917226 59.88%  github.com/go-chi/chi/v5.(*ChainHandler).ServeHTTP
         0     0% 64.63%  -44749955 61.02%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 64.63%  -43982763 59.97%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 64.63%   -5036804  6.87%  github.com/jackc/pgx/v5.(*Conn).Query
         0     0% 64.63%    -597594  0.81%  github.com/jackc/pgx/v5.(*Conn).getStatementDescription
         0     0% 64.63%   -3188391  4.35%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).Build
         0     0% 64.63%   -8331545 11.36%  github.com/jackc/pgx/v5.ConnectConfig
         0     0% 64.63%   -5762616  7.86%  github.com/jackc/pgx/v5.ParseConfig (inline)
         0     0% 64.63%    -731705     1%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecPrepared
         0     0% 64.63%    -519072  0.71%  github.com/jackc/pgx/v5/pgconn.(*PgConn).execExtendedSuffix
         0     0% 64.63%   -1875357  2.56%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
         0     0% 64.63%    -458758  0.63%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).readUntilRowDescription
         0     0% 64.63%    -657456   0.9%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).receiveMessage
         0     0% 64.63%    -450594  0.61%  github.com/jackc/pgx/v5/pgconn.(*scramClient).recvServerFinalMessage
         0     0% 64.63%   -7548232 10.29%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
         0     0% 64.63%    -942158  1.28%  github.com/jackc/pgx/v5/pgconn.computeHMAC
         0     0% 64.63%   -1875357  2.56%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
         0     0% 64.63%   -3051848  4.16%  github.com/jackc/pgx/v5/pgtype.(*Map).Encode
         0     0% 64.63%    -432829  0.59%  io.Copy (inline)
         0     0% 64.63%    -607595  0.83%  io.CopyN
         0     0% 64.63%    -327685  0.45%  io.ReadAtLeast
         0     0% 64.63%    -327685  0.45%  io.ReadFull (inline)
         0     0% 64.63%    -432829  0.59%  io.copyBuffer
         0     0% 64.63%    -432829  0.59%  io.discard.ReadFrom
         0     0% 64.63%   -1515747  2.07%  main.main
         0     0% 64.63%   -1515747  2.07%  main.run
         0     0% 64.63%    -393222  0.54%  net.(*TCPAddr).String
         0     0% 64.63%   -1099161  1.50%  net.(*TCPListener).Accept
         0     0% 64.63%   -1099161  1.50%  net.(*TCPListener).accept
         0     0% 64.63%    -902552  1.23%  net.(*netFD).accept
         0     0% 64.63%    -375298  0.51%  net.(*netFD).dial
         0     0% 64.63%    -562907  0.77%  net.(*sysDialer).dialParallel
         0     0% 64.63%    -562907  0.77%  net.(*sysDialer).dialSerial
         0     0% 64.63%    -510474   0.7%  net.(*sysDialer).dialTCP
         0     0% 64.63%    -510474   0.7%  net.(*sysDialer).doDialTCP (inline)
         0     0% 64.63%    -510474   0.7%  net.(*sysDialer).doDialTCPProto
         0     0% 64.63%    -444938  0.61%  net.internetSocket
         0     0% 64.63%    -444938  0.61%  net.socket
         0     0% 64.63%   -1545212  2.11%  net/http.(*Server).ListenAndServe
         0     0% 64.63%    -327685  0.45%  net/http.(*body).Read
         0     0% 64.63%    -327685  0.45%  net/http.(*body).readLocked
         0     0% 64.63%    -875840  1.19%  net/http.(*chunkWriter).Write
         0     0% 64.63%    -875840  1.19%  net/http.(*chunkWriter).writeHeader
         0     0% 64.63%    -507480  0.69%  net/http.(*response).WriteHeader
         0     0% 64.63%    -954630  1.30%  net/http.(*response).finishRequest
         0     0% 64.63%  -48850232 66.61%  net/http.HandlerFunc.ServeHTTP
         0     0% 64.63%    -327685  0.45%  net/http.Header.Get
         0     0% 64.63%    -759698  1.04%  net/http.Header.Set (inline)
         0     0% 64.63%   -1545212  2.11%  net/http.ListenAndServe (inline)
         0     0% 64.63%    -606190  0.83%  net/http.newBufioWriterSize
         0     0% 64.63%    -303928  0.41%  net/http.newTextprotoReader
         0     0% 64.63%  -48850232 66.61%  net/http.serverHandler.ServeHTTP
         0     0% 64.63%   -1409353  1.92%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0% 64.63%    -425990  0.58%  net/textproto.CanonicalMIMEHeaderKey
         0     0% 64.63%    -327685  0.45%  net/textproto.MIMEHeader.Get
         0     0% 64.63%   -1032224  1.41%  os.OpenFile
         0     0% 64.63%    -391570  0.53%  os.Stat
         0     0% 64.63%    -917530  1.25%  os.WriteFile
         0     0% 64.63%    -529758  0.72%  os.ignoringEINTR (inline)
         0     0% 64.63%    -376838  0.51%  os.open (inline)
         0     0% 64.63%    -376838  0.51%  os.openFileNolog.func1 (inline)
         0     0% 64.63%   -1306486  1.78%  reflect.Value.Grow
         0     0% 64.63%   -1306486  1.78%  reflect.Value.grow
         0     0% 64.63%    -737303  1.01%  strings.(*Builder).Grow
         0     0% 64.63%   -1963893  2.68%  sync.(*Pool).Get
         0     0% 64.63%    -559682  0.76%  sync.(*Pool).Put
         0     0% 64.63%   -1247073  1.70%  sync.(*Pool).pin
         0     0% 64.63%    -595295  0.81%  syscall.BytePtrFromString (inline)
         0     0% 64.63%    -376838  0.51%  syscall.Open
```