import json
import pandas as pd
from datetime import datetime
import matplotlib.pyplot as plt

# 获取当前时间字符串
current_time_str = datetime.now().strftime("%Y%m%d_%H%M%S")

# 读取并过滤数据
data = []
with open('/Users/clay/workspace/eigenda/tools/perf/eigenda_tx_retrieve/eigenda20240803162445_mainnet_20240803162445.txt', 'r') as file:
    i =0
    for line in file:
        i = i + 1
        print("i:", i)
        if not line.startswith("{"):
            continue
        log_entry = json.loads(line)
        data.append(log_entry)

# 将数据转化为DataFrame
confirm_time_counts = {}
send_request_time_counts = {}

includeConfirmData = []
for item in data:
    if item["confirm_time"] > 0:
        includeConfirmData.append(item)
    item["send_request_time"] = datetime.fromtimestamp(item["send_request_time"])
    item["receive_request_time"] = datetime.fromtimestamp(item["receive_request_time"])
    item["confirm_time"] = datetime.fromtimestamp(item["confirm_time"])

# 绘制每分钟完成处理的blob数量折线图
df = pd.DataFrame(includeConfirmData)
df.set_index('confirm_time', inplace=True)
confirm_time_counts = df.resample('1min').count()

plt.figure(figsize=(14, 7))
plt.plot(confirm_time_counts.index, confirm_time_counts.request_id, marker='o', linestyle='-', color='#87CEEB', label='Number of Completed Blobs Per Minute')
plt.title('Number of Completed Blobs Per Minute')
plt.xlabel('Time (Minute)')
plt.ylabel('Number of Blobs')
plt.grid(True)
plt.savefig(f'completed_blobs_per_minute_{current_time_str}.png')
plt.show()

# 绘制每分钟发送的请求数
df2 = pd.DataFrame(data)
df2.set_index('send_request_time', inplace=True)
send_request_time_counts = df2.resample('1min').count()

plt.figure(figsize=(14, 7))
plt.plot(send_request_time_counts.index, send_request_time_counts.request_id, marker='o', linestyle='-', color='#87CEEB', label='Number of Completed Blobs Per Minute')
plt.title('Number of sent Blobs Per Minute')
plt.xlabel('Time (Minute)')
plt.ylabel('Number of Blobs')
plt.grid(True)
plt.savefig(f'sent_blobs_per_minute_{current_time_str}.png')
plt.show()

# 绘制每分钟成功的请求数
df3 = pd.DataFrame(includeConfirmData)
df3.set_index('receive_request_time', inplace=True)
receive_request_time_counts = df3.resample('1min').count()

plt.figure(figsize=(14, 7))
plt.plot(receive_request_time_counts.index, receive_request_time_counts.request_id, marker='o', linestyle='-', color='#87CEEB', label='Number of Completed Blobs Per Minute')
plt.title('The Number of Sent Successfully Blobs Per Minute')
plt.xlabel('Time (Minute)')
plt.ylabel('Number of Blobs')
plt.grid(True)
plt.savefig(f'sent_successfully_blobs_per_minute_{current_time_str}.png')
plt.show()
