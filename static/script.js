document.addEventListener('DOMContentLoaded', function () {
    const loadingElement = document.getElementById('loading');
    const errorElement = document.getElementById('error');
    const gamesElement = document.getElementById('games');

    // 增加篩選表單元素，先隱藏
    const filterForm = document.createElement('div');
    filterForm.id = 'filterForm';
    filterForm.style.display = 'none'; // 初始隱藏篩選表單
    filterForm.innerHTML = `
        <label for="minOdds">最小受讓分:</label>
        <input type="number" id="minOdds" name="minOdds">
        <label for="maxOdds">最大受讓分:</label>
        <input type="number" id="maxOdds" name="maxOdds">
    `;
    document.body.insertBefore(filterForm, gamesElement);

    // 顯示加載狀態
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
            // 隱藏加載狀態
            loadingElement.style.display = 'none';
            renderGames(gamesData.matches);
            // 資料載入成功後，顯示篩選表單
            filterForm.style.display = 'block';
            setupFilterListeners(); // 設置篩選監聽
        })
        .catch(error => {
            console.error('Error fetching data: ', error);
            // 顯示錯誤信息
            errorElement.textContent = `Error: ${error.message}`;
            errorElement.style.display = 'block';
            // 隱藏加載狀態
            loadingElement.style.display = 'none';
        });
});

function renderGames(matches) {
    const gamesElement = document.getElementById('games');
    gamesElement.innerHTML = ''; // 清空現有內容
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
            <th>客隊傷兵名單</th>
            <th>主隊傷兵名單</th>
            <th>客隊近五場過盤</th>
            <th>主隊近五場過盤</th>
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
            <td>${formatInjuryList(match.awayInjuries)}</td>
            <td>${formatInjuryList(match.homeInjuries)}</td>
            <td>${formatDishResult(match.awayDish)}</td>
            <td>${formatDishResult(match.homeDish)}</td>
        `;
        row.dataset.currentOdds = match.currentOdds;
        tbody.appendChild(row);
    });
    table.appendChild(tbody);
    gamesElement.appendChild(table);
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

function formatInjuryList(injuries) {
    if (!injuries || injuries.length === 0) {
        return '無';
    }

    let injuryTable = '<table class="injury-table">';
    injuryTable += '<tr><th>球員</th><th>狀態</th></tr>'; // 表頭

    injuries.forEach(injury => {
        let playerNameLink = `<a href="${injury.link}" target="_blank">${injury.name}</a>`;
        let statusClass = injury.status === 'GTD' ? 'status-gtd' : 'status-out';
        injuryTable += `<tr><td>${playerNameLink}</td><td class="${statusClass}">${injury.status}</td></tr>`;
    });

    injuryTable += '</table>';
    return injuryTable;
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