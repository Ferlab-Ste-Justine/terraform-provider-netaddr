etcdctl \
--endpoints="127.0.0.1:32379" \
--cacert="certs/ca.pem" \
--key="certs/root.key" \
--cert="certs/root.pem" \
get --prefix "/test/mac/"