import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
    // 仮想ユーザー数を徐々に増やしていく
    vus: 50, // まずは50人から
    duration: '30s',
  };

const BASE_URL = 'http://localhost:8080';

export default function () {
  // 1. 商品一覧を取得
  const listRes = http.get(`${BASE_URL}/items`);
  check(listRes, { 'list status was 200': (r) => r.status == 200 });
  sleep(1);

  // 2. 特定の商品(ID=2)を取得
  const itemRes = http.get(`${BASE_URL}/items/2`);
  check(itemRes, { 'item status was 200': (r) => r.status == 200 });
  sleep(1);
  
  // 3. 新しい商品を登録
  const payload = JSON.stringify({ name: `New Item from k6 ${__VU}` });
  const params = { headers: { 'Content-Type': 'application/json' } };
  const createRes = http.post(`${BASE_URL}/items`, payload, params);
  check(createRes, { 'create status was 201': (r) => r.status == 201 });
  sleep(1);
}