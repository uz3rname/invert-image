const imageObjects = [];

const sleep = n => new Promise(r => setTimeout(r, n));

const showError = msg => {
  const elt = document.getElementById('status');
  elt.innerText = msg;
  elt.hidden = false;
  elt.style.color = '#a44';
};

const showInfo = msg => {
  const elt = document.getElementById('status');
  elt.innerText = msg;
  elt.hidden = false;
  elt.style.color = 'gray';
};

const clearFileList = () => {
  const elt = document.getElementById('last-image-list');
  elt.innerHTML = '';
  imageObjects.length = 0;
};

const openImage = data => {
  const newTab = window.open();
  newTab.document.body.innerHTML = `<img src="${data}">`;
};

const addFileToList = (obj, prepend) => {
  if (imageObjects.find(x => x.id === obj.id)) {
    return;
  }
  const elt = document.getElementById('last-image-list');

  const tr = document.createElement('tr');
  const origData =
    `data:${obj.original.mimeType};base64,${obj.original.base64}`;
  const negData =
    `data:${obj.negative.mimeType};base64,${obj.negative.base64}`;

  tr.innerHTML = `
    <td>${new Date(obj.createdAt).toLocaleString()}</td>
    <td><img style="cursor: hand" onclick="openImage('${origData}');" src="${origData}"></td>
    <td><img style="cursor: hand" onclick="openImage('${negData}');" src="${negData}"></td>`;
  if (prepend) {
    elt.prepend(tr);
    imageObjects.unshift(obj);
  } else {
    elt.appendChild(tr);
    imageObjects.push(obj);
  }
};

const sendImageData = () => {
  const elt = document.getElementById('file-input');
  
  if (elt.files.length === 0) {
    showError('No file selected!');
    return;
  }

  const reader = new FileReader();
  reader.onloadend = async () => {
    const array = new Uint8Array(reader.result);
    const data = Base64.fromUint8Array(array);
    elt.files[0];
    const resp = await fetch('/api/negative_image', {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
      },
      body: JSON.stringify({
        data: data.replace(/[=]+$/, ''),
      }),
    });
    const json = await resp.json();
    switch (json.status) {
      case 'ok':
        addFileToList(json.pair, true);
        showInfo('Success!');
        return;
      case 'defered':
        showInfo(
          `Image is too large, adding to queue, taskId=${json.taskId}`,
        );
        (async (taskId) => {
          while (true) {
            const resp = await fetch(`/api/get_task_status?taskId=${taskId}`);
            const json = await resp.json();
            if (json.taskStatus !== 'new' && json.taskStatus !== 'running') {
              if (json.taskStatus === 'done') {
                showInfo(`Task ${taskId} completed! Reloading`);
                loadLastImages();
                return;
              } else if (json.taskStatus === 'failed') {
                showError(`Task ${taskId} failed!`);
                return;
              } else if (json.taskStatus === 'canceled') {
                showInfo(`Task ${taskId} was cancelled`);
                return;
              }
            }
            await sleep(1000);
          }
        })(json.taskId);
        return;
      default:
        showError(json.message);
        return;
    }
  }
  reader.readAsArrayBuffer(elt.files[0]);
};

const loadLastImages = async () => {
  const resp = await fetch('/api/get_last_images');
  const json = await resp.json();
  if (json.status !== 'ok') {
    showError(json.message);
    return;
  }
  clearFileList();
  for (const item of json.items) {
    addFileToList(item);
  }
  showInfo('Updated!')
};

window.onload = () => {
  loadLastImages();
}
