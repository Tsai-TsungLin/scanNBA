document.addEventListener('DOMContentLoaded', function () {
    const loadingElement = document.getElementById('loading');
    const errorElement = document.getElementById('error');
    const gamesElement = document.getElementById('games');

    const filterForm = document.createElement('div');
    filterForm.id = 'filterForm';
    filterForm.style.display = 'none';
    filterForm.innerHTML = `
        <label for="minOdds">最小受讓分:</label>
        <input type="number" id="minOdds" name="minOdds">
        <label for="maxOdds">最大受讓分:</label>
        <input type="number" id="maxOdds" name="maxOdds">
    `;
    document.body.insertBefore(filterForm, gamesElement);

    loadingElement.style.display = 'block';
    errorElement.style.display = 'none';

    fetch('/api/games')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(gamesData => {
            loadingElement.style.display = 'none';
            renderGames(gamesData.matches);
            filterForm.style.display = 'block';
            setupFilterListeners();
        })
        .catch(error => {
            console.error('Error fetching data: ', error);
            errorElement.textContent = `Error: ${error.message}`;
            errorElement.style.display = 'block';
            loadingElement.style.display = 'none';
        });
});

function renderGames(matches) {
    const gamesElement = document.getElementById('games');
    gamesElement.innerHTML = '';
    const table = document.createElement('table');
    table.className = 'games-table';
    const thead = document.createElement('thead');
    thead.innerHTML = `
        <tr>
            <th>時間</th>
            <th>主隊</th>
            <th>客隊</th>
            <th>初盤受讓分</th>
            <th>初盤大小分</th>
            <th>目前受讓分</th>
            <th>目前大小分</th>
            <th>主隊大小分</th>
            <th>客隊大小分</th>
            <th>主隊傷兵名單</th>
            <th>客隊傷兵名單</th>
            <th>主隊近五場過盤</th>
            <th>客隊近五場過盤</th>
        </tr>`;
    table.appendChild(thead);
    const tbody = document.createElement('tbody');
    matches.forEach(match => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${match.time}</td>
            <td>${match.homeTeam}</td>
            <td>${match.awayTeam}</td>
            <td>${match.initialOdds}</td>
            <td>${match.initialOverUnder}</td>
            <td>${match.currentOdds}</td>
            <td>${match.currentOverUnder}</td>
            <td>${match.homeOverUnder}</td>
            <td>${match.awayOverUnder}</td>
            <td>${formatInjuryList(match.homeInjuries, 'awayInjuryList' + match.homeTeam)}</td>
            <td>${formatInjuryList(match.awayInjuries, 'homeInjuryList' + match.awayTeam)}</td>
            <td>${formatDishResult(match.homeDish)}</td>
            <td>${formatDishResult(match.awayDish)}</td>
        `;
        tbody.appendChild(row);
        row.dataset.currentOdds = match.currentOdds;
    });
    table.appendChild(tbody);
    gamesElement.appendChild(table);
}
// 更新 toggleInjuryList 函數以切換按鈕文字
function toggleInjuryList(listId) {
    var element = document.getElementById(listId);
    var button = document.getElementById('btn' + listId); // 根據listId獲取對應的按鈕
    if (element.style.display === "none") {
        element.style.display = "block";
        button.textContent = "收起"; // 當列表展開時，按鈕顯示“收起”
    } else {
        element.style.display = "none";
        button.textContent = "展開"; // 當列表收起時，按鈕顯示“展開”
    }
}

// 篩選功能啟動時的函數
function setupFilterListeners() {
    const minOddsInput = document.getElementById('minOdds');
    const maxOddsInput = document.getElementById('maxOdds');

    // 監聽輸入變化，自動觸發篩選
    minOddsInput.addEventListener('input', autoFilterGames);
    maxOddsInput.addEventListener('input', autoFilterGames);
}

// 自動篩選函數，當兩個輸入框都有值時觸發篩選
function autoFilterGames() {
    const minOdds = document.getElementById('minOdds').value;
    const maxOdds = document.getElementById('maxOdds').value;

    // 當兩個輸入框都有值或都是空白時，進行篩選
    if ((minOdds !== '' && maxOdds !== '') || (minOdds === '' && maxOdds === '')) {
        filterGames();
    }
}

// 篩選函數
function filterGames() {
    const minOddsValue = document.getElementById('minOdds').value;
    const maxOddsValue = document.getElementById('maxOdds').value;
    const rows = document.querySelectorAll('#games tbody tr');
    rows.forEach(row => {
        if (row.dataset.currentOdds) {
            const currentOdds = parseFloat(row.dataset.currentOdds);
            const displayRow = (!minOddsValue && !maxOddsValue) || // 兩個輸入框都是空的
                (minOddsValue !== '' && maxOddsValue !== '' && currentOdds >= parseFloat(minOddsValue) && currentOdds <= parseFloat(maxOddsValue)); // 兩個輸入框都有值且符合條件

            row.style.display = displayRow ? '' : 'none'; // 根據條件顯示或隱藏行
        }
    });
}

function formatInjuryList(injuries, listId) {
    // 檢查是否有傷兵數據
    if (!injuries || injuries.length === 0) {
        // 返回“無傷兵名單”並置中顯示，並賦予特定的class
        return '<div class="no-injuries">無傷兵名單</div>';
    }

    let buttonHTML = `<div class="injury-list-button-container"><button id="btn${listId}" class="injury-list-toggle-button" onclick="toggleInjuryList('${listId}')">展開</button></div>`;
    let injuryTableHTML = `<div id="${listId}" class="injury-list" style="display: none;">`;
    injuryTableHTML += '<table class="injury-table">';
    injuryTableHTML += '<tr><th>球員</th><th>狀態</th></tr>';

    injuries.forEach(injury => {
        let playerNameLink = `<a href="${injury.link}" target="_blank">${injury.name}</a>`;
        let statusClass = injury.status === 'GTD' ? 'status-gtd' : 'status-out';
        injuryTableHTML += `<tr><td>${playerNameLink}</td><td class="${statusClass}">${injury.status}</td></tr>`;
    });

    injuryTableHTML += '</table></div>';

    // 返回按鈕置中的HTML，放在傷兵名單的上方
    return buttonHTML + injuryTableHTML;
}


function formatDishResult(dishResults) {
    // 检查 dishResults 是否为数组
    if (!Array.isArray(dishResults)) {
        return '';
    }

    return dishResults.map(dish => {
        dish = dish.trim(); // 清理空格
        if (dish === '贏') {
            return '<span style="color: green;">贏</span>';
        } else if (dish === '輸') {
            return '<span style="color: red;">輸</span>';
        }
        return dish;
    }).join(', ');
}