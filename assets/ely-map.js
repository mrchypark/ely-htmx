(() => {
  const ns = 'http://www.w3.org/2000/svg';
  const svg = (name, attrs = {}) => { const node = document.createElementNS(ns, name); Object.entries(attrs).forEach(([key, value]) => node.setAttribute(key, value)); return node; };
  const point = (lat, lon) => [((Number(lon) + 180) / 360) * 100, ((90 - Number(lat)) / 180) * 100];
  const coordinate = (lat, lon) => `${Number(lat).toFixed(4)}, ${Number(lon).toFixed(4)}`;
  const parseCoordinate = value => { const match = String(value || '').match(/(-?\d+(?:\.\d+)?)\s*,\s*(-?\d+(?:\.\d+)?)/); return match ? [Number(match[1]), Number(match[2])] : null; };
  const pathFor = geometry => { const polygons = geometry?.type === 'Polygon' ? [geometry.coordinates] : geometry?.type === 'MultiPolygon' ? geometry.coordinates : []; return polygons.map(polygon => polygon.map(ring => `M${ring.map(([lon, lat]) => point(lat, lon).join(' ')).join(' ')}Z`).join('')).join(''); };
  const localFeatures = async () => { const response = await fetch('/assets/ely/world-land.json', { credentials: 'same-origin' }); if (!response.ok) throw new Error('world map unavailable'); const data = await response.json(); if (data.provenance?.source !== 'Natural Earth' || !Array.isArray(data.features)) throw new Error('invalid world map'); return data.features; };
  const init = root => {
    if (root.dataset.mapReady) return; root.dataset.mapReady = 'true';
    const mode = root.dataset.mapMode || 'markers', canvas = root.querySelector('[data-map-canvas]'), viewport = root.querySelector('[data-map-viewport]'), markers = root.querySelector('[data-map-markers]'), popup = root.querySelector('[data-map-popup]'), output = root.querySelector('[data-map-coordinate]'), input = root.querySelector('[data-map-value]'), status = root.querySelector('[data-map-status]'), crosshair = root.querySelector('.ely-map__crosshair');
    let scale = 1, pan = [0, 0], drag, popupTrigger, selected = parseCoordinate(input?.value);
    const apply = () => { viewport.setAttribute('transform', `translate(${pan[0]} ${pan[1]}) scale(${scale})`); markers.querySelectorAll('[data-map-point]').forEach(marker => { const [x, y] = point(marker.dataset.lat, marker.dataset.lon); marker.style.left = `${x * scale + pan[0]}%`; marker.style.top = `${y * scale + pan[1]}%`; }); if (crosshair && selected) { const [x, y] = point(...selected); crosshair.style.left = `${x * scale + pan[0]}%`; crosshair.style.top = `${y * scale + pan[1]}%`; } };
    const setCoordinate = (lat, lon) => { selected = [Math.max(-90, Math.min(90, lat)), ((lon + 540) % 360) - 180]; const value = coordinate(...selected); if (input) { input.value = value; input.dispatchEvent(new Event('input', { bubbles: true })); } if (output) output.value = value; apply(); };
    const drawFeatures = features => { const layer = root.querySelector('[data-map-geojson-layer]'), land = root.querySelector('[data-map-land]'); const values = new Map([...markers.querySelectorAll('[data-map-point]')].map(marker => [String(marker.getAttribute('aria-label')).toLowerCase(), Number(marker.dataset.value) || 0])); const max = Math.max(1, ...values.values()); features.forEach(feature => { const country = svg('path', { d: pathFor(feature.geometry) }); const value = values.get(String(feature.properties?.name || feature.properties?.iso_a2 || '').toLowerCase()) || 0; if (mode === 'choropleth') { country.style.setProperty('--ely-map-level', String(Math.round(value / max * 6))); layer.append(country); } else land.append(country); }); };
    const source = root.querySelector('template[data-map-geojson]'); const inline = source?.content.textContent.trim(); Promise.resolve().then(() => inline ? JSON.parse(inline).features || [] : localFeatures()).then(drawFeatures).catch(() => { root.dataset.mapGeometry = 'unavailable'; if (status) { status.hidden = false; status.textContent = 'Map geometry is unavailable.'; } });
    const routes = root.querySelector('[data-map-routes]'); root.querySelectorAll('[data-map-route]').forEach(route => { const stops = [...route.querySelectorAll('[data-map-point]')]; if (stops.length > 1) routes.append(svg('path', { d: stops.map((stop, index) => `${index ? 'L' : 'M'}${point(stop.dataset.lat, stop.dataset.lon).join(' ')}`).join(' ') })); });
    if (mode === 'heatmap') markers.querySelectorAll('[data-map-point]').forEach(marker => { const [cx, cy] = point(marker.dataset.lat, marker.dataset.lon); root.querySelector('[data-map-heat]').append(svg('circle', { cx, cy, r: Math.max(2, Number(marker.dataset.value) || 5) / 2 })); marker.hidden = true; });
    if (mode === 'clusters') { const groups = new Map(); markers.querySelectorAll('[data-map-point]').forEach(marker => { const [x, y] = point(marker.dataset.lat, marker.dataset.lon); const key = `${Math.round(x / 8)}:${Math.round(y / 8)}`; groups.set(key, [...(groups.get(key) || []), marker]); }); groups.forEach(group => { if (group.length < 2) return; group.slice(1).forEach(marker => marker.hidden = true); group[0].dataset.cluster = group.length; group[0].textContent = group.length; group[0].setAttribute('aria-label', `${group.length} locations`); }); }
    const closePopup = () => { popup.hidden = true; popupTrigger?.focus(); popupTrigger = null; };
    const showPopup = marker => { popupTrigger = marker; popup.hidden = false; popup.querySelector('[data-map-popup-label]').textContent = marker.getAttribute('aria-label'); popup.querySelector('[data-map-popup-value]').textContent = marker.dataset.value; popup.style.left = marker.style.left; popup.style.top = marker.style.top; popup.querySelector('[data-map-close]').focus(); };
    markers.addEventListener('click', event => { const marker = event.target.closest('[data-map-point]'); if (marker && !root.dataset.mapDisabled) showPopup(marker); }); root.querySelector('[data-map-close]').addEventListener('click', closePopup);
    root.querySelector('[data-map-zoom="in"]').addEventListener('click', () => { scale = Math.min(4, scale * 1.25); apply(); }); root.querySelector('[data-map-zoom="out"]').addEventListener('click', () => { scale = Math.max(1, scale / 1.25); apply(); }); root.querySelector('[data-map-zoom="reset"]').addEventListener('click', () => { scale = 1; pan = [0, 0]; apply(); });
    const stopDrag = () => { drag = null; canvas.classList.remove('is-dragging'); };
    canvas.addEventListener('pointerdown', event => { if (root.dataset.mapDisabled || event.target.closest('[data-map-point],[data-map-popup]')) return; drag = [event.clientX, event.clientY, ...pan]; canvas.setPointerCapture(event.pointerId); canvas.classList.add('is-dragging'); }); canvas.addEventListener('pointermove', event => { if (!drag) return; const box = canvas.getBoundingClientRect(); pan = [drag[2] + (event.clientX - drag[0]) / box.width * 100, drag[3] + (event.clientY - drag[1]) / box.height * 100]; apply(); }); canvas.addEventListener('pointerup', stopDrag); canvas.addEventListener('pointercancel', stopDrag); canvas.addEventListener('lostpointercapture', stopDrag);
    canvas.addEventListener('click', event => { if (mode !== 'location' || root.dataset.mapDisabled || event.target.closest('[data-map-point],[data-map-popup]')) return; const box = canvas.getBoundingClientRect(), x = ((event.clientX - box.left) / box.width * 100 - pan[0]) / scale, y = ((event.clientY - box.top) / box.height * 100 - pan[1]) / scale; setCoordinate(90 - y / 100 * 180, x / 100 * 360 - 180); }); canvas.addEventListener('keydown', event => {
      if (root.dataset.mapDisabled || event.target !== canvas) return;
      if (['+', '=', '-', '0'].includes(event.key)) {
        event.preventDefault();
        scale = event.key === '-' ? Math.max(1, scale / 1.25) : event.key === '0' ? 1 : Math.min(4, scale * 1.25);
        if (event.key === '0') pan = [0, 0];
        apply(); return;
      }
      if (!['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(event.key)) return;
      event.preventDefault();
      if (mode !== 'location') {
        pan[0] += event.key === 'ArrowLeft' ? 5 : event.key === 'ArrowRight' ? -5 : 0;
        pan[1] += event.key === 'ArrowUp' ? 5 : event.key === 'ArrowDown' ? -5 : 0;
        apply(); return;
      }
      const current = selected || [0, 0], delta = event.shiftKey ? 1 : .1;
      setCoordinate(current[0] + (event.key === 'ArrowUp' ? delta : event.key === 'ArrowDown' ? -delta : 0), current[1] + (event.key === 'ArrowRight' ? delta : event.key === 'ArrowLeft' ? -delta : 0));
    }); document.addEventListener('keydown', event => { if (event.key === 'Escape' && !popup.hidden && root.contains(document.activeElement)) closePopup(); }); if (selected) setCoordinate(...selected); apply();
  };
  const start = () => document.querySelectorAll('[data-ui-component="map-view"]').forEach(init); start(); document.addEventListener('htmx:after:swap', start);
})();
