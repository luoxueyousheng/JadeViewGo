'use strict';
const $ = s => document.querySelector(s);
const hasJade = typeof window.jade !== 'undefined';

/* ============ Toast（DESIGN.md §12） ============ */
const TOAST_ICON = { info: '#i-info', success: '#i-check', warning: '#i-warn', error: '#i-error' };
const TOAST_DUR  = { info: 4000, success: 4000, warning: 6000, error: 0 };
const toastById = new Map();

function showToast({ level = 'info', title = '', message = '', duration, id, action } = {}) {
  const box = $('#toasts');
  if (id && toastById.has(id)) { toastById.get(id).remove(); toastById.delete(id); }
  while (box.children.length >= 4) box.firstChild.remove();   // 上限，挤掉最旧

  const el = document.createElement('div');
  el.className = 'toast ' + level;
  if (level === 'error') el.setAttribute('aria-live', 'assertive');
  el.innerHTML =
    `<svg class="icon toast-icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><use href="${TOAST_ICON[level] || TOAST_ICON.info}"/></svg>` +
    `<div class="body">${title ? `<div class="title"></div>` : ''}<div class="msg"></div></div>` +
    (action ? `<button class="btn" style="height:24px;padding:0 8px;font-size:12px"></button>` : '') +
    `<button class="close" aria-label="关闭"><svg class="icon" width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><use href="#i-close"/></svg></button>`;
  if (title) el.querySelector('.title').textContent = title;
  el.querySelector('.msg').textContent = message;
  if (action) {
    const b = el.querySelector('.body + .btn');
    b.textContent = action.label;
    b.onclick = () => { if (hasJade) jade.invoke(action.command).catch(() => {}); dismiss(); };
  }

  let ms = duration ?? TOAST_DUR[level] ?? 4000;
  if (action) ms = 0;                                          // 含操作的不自动消失
  let timer = null, left = ms, started = 0;
  const dismiss = () => {
    if (timer) clearTimeout(timer);
    el.classList.add('toast-out');
    el.addEventListener('animationend', () => el.remove(), { once: true });
    if (id) toastById.delete(id);
  };
  const arm = () => { if (left > 0) { started = Date.now(); timer = setTimeout(dismiss, left); } };
  el.addEventListener('mouseenter', () => { if (timer) { clearTimeout(timer); timer = null; left -= Date.now() - started; } });
  el.addEventListener('mouseleave', arm);
  el.querySelector('.close').onclick = dismiss;

  box.appendChild(el);
  if (id) toastById.set(id, el);
  arm();
}

/* ============ invoke 封装 ============ */
async function inv(channel, payload = {}, silent = false) {
  if (!hasJade) { showToast({ level: 'error', title: '环境错误', message: 'jade 对象不可用（不在 JadeView 内运行）。' }); return null; }
  const t0 = performance.now();
  try {
    const res = await jade.invoke(channel, payload, { timeout: 8000 });
    const ms = Math.round(performance.now() - t0);
    ipcLog(`invoke('${channel}') ${ms}ms → ${typeof res === 'string' ? res : JSON.stringify(res)}`, true);
    if (!silent) showToast({ level: 'success', title: channel, message: String(res) });
    return res;
  } catch (e) {
    ipcLog(`invoke('${channel}') 失败: ${e}`);
    showToast({ level: 'error', title: channel + ' 失败', message: String(e) });
    return null;
  }
}

function ipcLog(msg, ok = false) {
  const el = $('#ipcLog');
  const t = new Date().toLocaleTimeString();
  const line = document.createElement('div');
  line.innerHTML = `<span class="t">[${t}]</span> <span class="${ok ? 'ok' : ''}"></span>`;
  line.lastElementChild.textContent = msg;
  el.appendChild(line);
  el.scrollTop = el.scrollHeight;
}

/* ============ NavigationView：共享指示条滑动（§11.7.1） ============ */
const nav = $('#nav'), navInd = $('#navInd');
function moveIndicator(item) {
  if (!item) return;
  const navRect = nav.getBoundingClientRect(), r = item.getBoundingClientRect();
  const barH = Math.max(12, r.height - 16);
  navInd.style.height = barH + 'px';
  navInd.style.transform = `translateY(${(r.top - navRect.top) + (r.height - barH) / 2}px)`;
}
nav.addEventListener('click', e => {
  const item = e.target.closest('.nav-item[data-page]');
  if (!item) return;
  nav.querySelectorAll('.nav-item[data-page]').forEach(x => {
    x.classList.toggle('active', x === item);
    x.setAttribute('aria-selected', x === item ? 'true' : 'false');
  });
  moveIndicator(item);
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  const page = $('#page-' + item.dataset.page);
  page.classList.add('active');
  segRefreshers.forEach(f => f());              // 页面变可见后重定位分段药丸
  if (item.dataset.page === 'overview') loadOverview();
});
$('#hamburger').onclick = () => {
  nav.toggleAttribute('data-collapsed');
  setTimeout(() => moveIndicator(nav.querySelector('.nav-item.active')), 260); // 过渡结束后重定位
};
new ResizeObserver(() => moveIndicator(nav.querySelector('.nav-item.active'))).observe(nav);

/* ============ 分段控件：共享药丸滑动 + 方向键（§6.8 B） ============ */
const segRefreshers = [];   // 页面从隐藏变可见后重定位药丸（display:none 时测不到宽度）
function initSeg(root, onChange) {
  const ind = root.querySelector('.seg-ind');
  const tabs = () => [...root.querySelectorAll('button[role="tab"]')];
  function select(btn, fire = true) {
    tabs().forEach(b => b.setAttribute('aria-selected', b === btn ? 'true' : 'false'));
    ind.style.width = btn.offsetWidth + 'px';
    ind.style.transform = `translateX(${btn.offsetLeft}px)`;   // offsetLeft 原点与 left:0 一致
    if (fire) onChange(btn.dataset.v);
  }
  root.addEventListener('click', e => { const b = e.target.closest('button[role="tab"]'); if (b) select(b); });
  root.addEventListener('keydown', e => {
    const list = tabs(), i = list.findIndex(b => b.getAttribute('aria-selected') === 'true');
    const next = { ArrowRight: i + 1, ArrowLeft: i - 1, Home: 0, End: list.length - 1 }[e.key];
    if (next == null) return;
    e.preventDefault();
    const b = list[(next + list.length) % list.length];
    select(b); b.focus();
  });
  const refresh = () => {
    const cur = root.querySelector('button[aria-selected="true"]') || tabs()[0];
    if (cur && root.offsetWidth > 0) select(cur, false);
  };
  segRefreshers.push(refresh);
  requestAnimationFrame(refresh);
  return { select };
}

/* ============ 运行环境（宿主经 "env" 通道下发，决定标题栏/材质适配） ============ */
const ENV = { os: 'windows', arch: '', win11: true };   // 独立预览/旧宿主：按 Windows 全功能展示
function applyPlatform() {
  document.documentElement.dataset.os = ENV.os;         // CSS 据此隐藏自绘标题栏（Linux 用系统标题栏）
  const hasBackdrop = ENV.os === 'windows' && ENV.win11;
  if (hasBackdrop) return;
  // 无 DWM 材质的平台：默认纯色背景，材质选项只留「纯色」
  currentBackdrop = 'none';
  document.querySelectorAll('#matList .mat').forEach(b => {
    if (b.dataset.v !== 'none') { b.disabled = true; b.setAttribute('aria-disabled', 'true'); }
  });
  const hint = document.querySelector('#matList + .hint');
  if (hint) hint.textContent = ENV.os === 'windows'
    ? '当前系统非 Windows 11，DWM 材质不可用，已改用纯色背景。'
    : `当前平台 ${ENV.os}：无 DWM 材质，已改用纯色背景（set_window_background_color）。`;
  const desc = $('#winDesc');
  if (desc) desc.textContent =
    `系统边框与标题栏（${ENV.os}）+ 纯色背景。DWM 材质与标题栏覆盖层为 Windows 专属，相关选项已停用。`;
}

/* ============ 主题（颜色模式） ============ */
const mqDark = matchMedia('(prefers-color-scheme: dark)');
let themeMode = 'system';

function effectiveDark() { return themeMode === 'dark' || (themeMode === 'system' && mqDark.matches); }
async function applyTheme() {
  const dark = effectiveDark();
  document.documentElement.dataset.theme = dark ? 'dark' : 'light';
  if (!hasJade) return;
  const mode = { light: 'Light', dark: 'Dark', system: 'System' }[themeMode];
  await inv('set-theme', { mode }, true);
  if (ENV.os === 'windows') await inv('apply-titlebar', { dark }, true);  // 标题栏覆盖层仅 Windows
  if (currentBackdrop === 'none') applyBackdrop('none', true);  // 纯色底随主题换色
}
mqDark.addEventListener('change', () => { if (themeMode === 'system') applyTheme(); });
if (hasJade) jade.on('theme-changed', () => { if (themeMode === 'system') applyTheme(); });

initSeg($('#segTheme'), v => { themeMode = v; applyTheme(); });

/* ============ 材质切换 ============ */
let currentBackdrop = 'mica';
async function applyBackdrop(type, silent = false) {
  currentBackdrop = type;
  document.documentElement.dataset.backdrop = type;
  document.querySelectorAll('#matList .mat')
    .forEach(b => b.setAttribute('aria-pressed', b.dataset.v === type ? 'true' : 'false'));
  const payload = { type };
  if (type === 'none') payload.color = effectiveDark() ? '#202020FF' : '#F3F3F3FF';
  await inv('set-backdrop', payload, silent);
}
$('#matList').addEventListener('click', e => {
  const b = e.target.closest('.mat');
  if (b) applyBackdrop(b.dataset.v);
});

/* ============ 缩放 ============ */
initSeg($('#segZoom'), v => inv('zoom', { level: parseFloat(v) }));

/* ============ IPC 测试页 ============ */
$('#btnInvoke').onclick = async () => {
  let payload = {};
  try { payload = JSON.parse($('#ipcPayload').value || '{}'); }
  catch { showToast({ level: 'warning', title: 'payload 无效', message: '请输入合法 JSON。' }); return; }
  const t0 = performance.now();
  const res = await inv($('#ipcChannel').value.trim(), payload, true);
  $('#lastLat').textContent = res == null ? '失败' : Math.round(performance.now() - t0) + ' ms';
};
$('#btnClearLog').onclick = () => $('#ipcLog').replaceChildren();
if (hasJade) {
  jade.on('push-demo', p => ipcLog(`← 宿主推送 push-demo: ${typeof p === 'string' ? p : JSON.stringify(p)}`, true));
  jade.on('toast', p => showToast(typeof p === 'string' ? JSON.parse(p) : p));
  jade.on('dialog-result', p => {
    const s = typeof p === 'string' ? p : JSON.stringify(p);
    ipcLog(`← 宿主推送 dialog-result: ${s}`, true);
    showToast({ level: 'info', title: '对话框结果', message: s });
  });
}

/* ============ 通用按钮绑定 ============ */
document.querySelectorAll('[data-inv]').forEach(b => {
  b.addEventListener('click', async () => {
    const res = await inv(b.dataset.inv, {}, true);
    if (res != null) {
      const out = { 'yaml-all': 1, 'yaml-get': 1, 'yaml-set': 1 }[b.dataset.inv] ? $('#yamlOut') : null;
      if (out) { out.textContent = String(res); }
      else showToast({ level: 'info', title: b.textContent.trim(), message: String(res) });
    }
  });
});
document.querySelectorAll('[data-toast]').forEach(b =>
  b.addEventListener('click', () => inv('demo-toast', { level: b.dataset.toast }, true)));
$('#swTop').addEventListener('change', () => inv('toggle-top', {}, false));

/* ============ 概览页数据 ============ */
async function loadOverview() {
  const info = await inv('sysinfo', {}, true);
  if (info) {
    $('#sysinfoKv').innerHTML = info.split('|').map(seg => {
      const [k, ...v] = seg.trim().split(/\s+/);
      return `<dt></dt><dd></dd>`;
    }).join('');
    const dts = $('#sysinfoKv').querySelectorAll('dt'), dds = $('#sysinfoKv').querySelectorAll('dd');
    info.split('|').forEach((seg, i) => {
      const s = seg.trim(), sp = s.indexOf(' ');
      dts[i].textContent = sp > 0 ? s.slice(0, sp) : s;
      dds[i].textContent = sp > 0 ? s.slice(sp + 1) : '';
    });
  }
  const disp = await inv('displays', {}, true);
  if (disp) { $('#displayKv').innerHTML = '<dt>信息</dt><dd></dd>'; $('#displayKv').querySelector('dd').textContent = disp; }
}

/* ============ 非激活态（窗口失焦变暗） ============ */
addEventListener('blur',  () => document.documentElement.dataset.inactive = '');
addEventListener('focus', () => delete document.documentElement.dataset.inactive);

/* ============ 启动 ============ */
moveIndicator(nav.querySelector('.nav-item.active'));
if (hasJade) {
  (async () => {
    const env = await inv('env', {}, true);           // 先取平台再适配（失败则按默认 Windows 处理）
    try { Object.assign(ENV, JSON.parse(env)); } catch { }
    applyPlatform();
    await applyTheme();             // System 模式：探测明暗；非 Win11 时顺带铺纯色背景
    loadOverview();
    showToast({ level: 'success', title: '已就绪', message: 'jade 对象可用，IPC 通道已连通。' });
  })();
} else {
  $('#titleSub').textContent = '独立预览（jade 不可用）';
  showToast({ level: 'warning', title: '独立预览', message: '未在 JadeView 内运行，仅可预览界面。', duration: 0 });
}
