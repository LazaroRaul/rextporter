sleep 5
if curl rextporter:8080 | grep -q 'go_memstats_frees_total counter'; then
  echo "Tests passed!"
  exit 1
else
  echo "Tests failed!"
  exit 0
fi