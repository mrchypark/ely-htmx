(() => {
  const status=(root,message)=>{ let out=root.querySelector('[data-mm-status]');if(!out){out=document.createElement('output');out.dataset.mmStatus='';root.append(out)}out.textContent=message; };
  document.addEventListener('click',event=>{
    const b=event.target.closest('button');const root=b?.closest('[data-mm]');if(!root)return;
    const action=b.dataset.mmAction;
    const image=root.querySelector('[data-mm-image]');
    if(action==='reset'&&root.querySelector('[data-mm-action="record"]')){const record=root.querySelector('[data-mm-action="record"]');record.setAttribute('aria-pressed','false');record.textContent='Start simulation';status(root,'00:00 · no screen capture');return;}
    if(action==='zoom'||action==='rotate'||action==='reset'){
      let scale=Number(root.dataset.mmScale)||1, rotation=Number(root.dataset.mmRotation)||0;
      if(action==='zoom')scale=scale>=2?1:scale+.25;
      if(action==='rotate')rotation=(rotation+90)%360;
      if(action==='reset'){scale=1;rotation=0}
      root.dataset.mmScale=scale;root.dataset.mmRotation=rotation;
      if(image)image.style.transform=`scale(${scale}) rotate(${rotation}deg)`;
      const model=root.querySelector('.mm-cube');if(model&&action==='reset')model.style.transform='';
      status(root,`${Math.round(scale*100)}% · ${rotation}°`);
    }
    if(action==='play'||action==='record'||action==='mute'){
      const pressed=b.getAttribute('aria-pressed')!=='true';b.setAttribute('aria-pressed',String(pressed));
      b.textContent=action==='mute'?(pressed?'Unmute':'Mute'):action==='record'?(pressed?'Stop simulation':'Start simulation'):(pressed?'Pause simulation':'Play simulation');
      status(root,pressed?'Local simulation active. No media device or file is opened.':'Simulation paused.');
    }
    if(action==='apply-crop')status(root,'Crop preview applied locally. Source image is unchanged.');
    if(action==='annotate'){const layer=root.querySelector('[data-mm-annotations]');if(layer&&!layer.firstChild){const note=document.createElement('span');note.className='mm-annotation';note.contentEditable='true';note.setAttribute('role','textbox');note.setAttribute('aria-label','Image annotation');note.textContent='Add a local note';layer.append(note)}}
    if(action==='clear-annotations')root.querySelector('[data-mm-annotations]')?.replaceChildren();
    if(action==='spectrum')root.querySelector('.mm-spectrum')?.classList.toggle('is-alternate');
    if(action==='level'){const meter=root.querySelector('meter');meter.value=meter.value===63?28:63}
    if(action==='camera'&&image)image.src=image.src.includes('atrium-olive')?'/assets/ely/assets/atrium-stair.jpg':'/assets/ely/assets/atrium-olive.jpg';
    if(action==='rotate-model'){const model=root.querySelector('.mm-cube');root.dataset.mmRotation=(Number(root.dataset.mmRotation)||0)+30;model.style.transform=`rotate(${root.dataset.mmRotation}deg)`}
    if(action==='unpin'){b.closest('article').remove();status(root,'Sample message unpinned.');}
    if(action==='reaction'){const selected=b.getAttribute('aria-pressed')!=='true';b.setAttribute('aria-pressed',String(selected));b.querySelector('span').textContent=selected?'3':'2'}
    if(action==='call-end'){b.disabled=true;status(root,'Local call simulation ended.');}
    if(b.dataset.mmTrack)status(root,`Selected: ${b.dataset.mmTrack} · simulated playback`);
    if(b.dataset.mmFolder){root.querySelectorAll('[data-mm-folder]').forEach(x=>x.removeAttribute('aria-current'));b.setAttribute('aria-current','page');status(root,b.dataset.mmFolder)}
    if(b.dataset.mmSticker)status(root,`Selected sample: ${b.dataset.mmSticker}`);
    if(b.dataset.mmThumbnail&&image)image.src=`/assets/ely/assets/${b.dataset.mmThumbnail}`;
  });
  document.addEventListener('input',event=>{
    const i=event.target,root=i.closest('[data-mm]');if(!root)return;
    if(i.matches('[data-mm-mail-filter]'))root.querySelectorAll('.mm-mail-list li').forEach(row=>{row.hidden=!row.textContent.toLowerCase().includes(i.value.toLowerCase())});
    if(i.matches('[data-mm-signature]'))root.querySelector('[data-mm-signature-preview]').textContent=i.value;
    if(i.matches('[data-mm-crop]'))root.querySelector('[data-mm-image]').style.width=`${i.value}%`;
  });
  document.addEventListener('change',event=>{
    const i=event.target,root=i.closest('[data-mm]');if(!root)return;
    if(i.matches('[data-mm-file]'))status(root,i.files?.[0]?`${i.files[0].name} · selected locally; not uploaded.`:'No file selected.');
    if(i.matches('[data-mm-aspect]'))root.querySelector('[data-mm-image]').style.aspectRatio=i.value;
    if(i.matches('[data-mm-snooze]'))status(root,`${i.value} · local reminder preview`);
  });
  document.addEventListener('submit',event=>{
    const form=event.target;if(!form.matches('[data-mm-compose]'))return;event.preventDefault();
    const i=form.querySelector('input'),host=form.closest('[data-mm]').querySelector('[data-mm-messages]');
    if(i.value.trim()){const p=document.createElement('p');p.className='demo-card';p.textContent=`You: ${i.value.trim()}`;host.append(p);i.value=''}
  });
})();
