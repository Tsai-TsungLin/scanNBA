document.addEventListener('DOMContentLoaded', function () {
    const loadingElement = document.getElementById('loading');
    const errorElement = document.getElementById('error');
    const gamesElement = document.getElementById('games');

    // 显示加载状态
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
            // 隐藏加载状态
            loadingElement.style.display = 'none';

            const gamesElement = document.getElementById('games');

            // 創建表格
            const table = document.createElement('table');
            table.className = 'games-table';

            // 創建表頭
            const thead = document.createElement('thead');
            thead.innerHTML = `
                <tr>
                    <th>時間</th>
                    <th>客隊</th>
                    <th>主隊</th>
                    <th>初盤受讓分</th>
                    <th>初盤大小分</th>
                    <th>目前受讓分</th>
                    <th>目前大小分</th>
                    <th>客隊傷兵名單</th>
                    <th>主隊傷兵名單</th>
                    <th>客隊近五場過盤</th>
                    <th>主隊近五場過盤</th>
                </tr>`;
            table.appendChild(thead);

            // 創建表身
            const tbody = document.createElement('tbody');

            gamesData.matches.forEach(match => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${match.time}</td>
                    <td>${match.awayTeam}</td>
                    <td>${match.homeTeam}</td>
                    <td>${match.initialOdds}</td>
                    <td>${match.initialOverUnder}</td>
                    <td>${match.currentOdds}</td>
                    <td>${match.currentOverUnder}</td>
                    <td>${formatInjuryList(match.awayInjuries)}</td>
                    <td>${formatInjuryList(match.homeInjuries)}</td>
                    <td>${formatDishResult(match.awayDish)}</td>
                    <td>${formatDishResult(match.homeDish)}</td>
                `;
                tbody.appendChild(row);
            });

            table.appendChild(tbody);
            gamesElement.appendChild(table);
        })
        .catch(error => {
            console.error('Error fetching data: ', error);
            // 显示错误信息
            errorElement.textContent = `Error: ${error.message}`;
            errorElement.style.display = 'block';
            // 隐藏加载状态
            loadingElement.style.display = 'none';
        });
});


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
