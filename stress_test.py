import socket
import threading
import time
import random

# 配置：目标是你本地的 Go 程序监听端口
TARGET_IP = '127.0.0.1'
TARGET_PORT = 8080

# 配置：并发线程数 (模拟多少个用户同时连接)
CONCURRENCY = 100 
# 配置：每个用户发多少次包
REQUESTS_PER_THREAD = 5

def attack(thread_id):
    """每个线程模拟一个用户，建立连接并发送数据"""
    try:
        # 1. 建立连接 (Connect)
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(5) # 5秒超时
        s.connect((TARGET_IP, TARGET_PORT))
        
        print(f"[Thread-{thread_id}] ✅ Connected!")

        # 2. 发送一些 HTTP 垃圾数据 (Send)
        # 我们不管百度回什么，只要发出去就算成功
        payload = f"GET / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: StressTest\r\n\r\n"
        
        for i in range(REQUESTS_PER_THREAD):
            s.sendall(payload.encode())
            # 接收一点数据，证明通路是活的
            data = s.recv(1024)
            if not data:
                break
            # 随机休眠一下，模拟真实人类行为
            time.sleep(random.uniform(0.1, 0.5))
            
        s.close()
        print(f"[Thread-{thread_id}] 👋 Closed normally")
        
    except Exception as e:
        print(f"[Thread-{thread_id}] ❌ Error: {e}")

def start_stress_test():
    threads = []
    print(f"🚀 Starting Stress Test: {CONCURRENCY} threads connecting to {TARGET_IP}:{TARGET_PORT}")
    
    start_time = time.time()

    # 启动所有线程
    for i in range(CONCURRENCY):
        t = threading.Thread(target=attack, args=(i,))
        threads.append(t)
        t.start()
        # 稍微错开一点启动时间，避免瞬间把 OS 端口耗尽
        time.sleep(0.01)

    # 等待所有线程结束
    for t in threads:
        t.join()

    end_time = time.time()
    print(f"\n🏁 Test Finished in {end_time - start_time:.2f} seconds!")

if __name__ == "__main__":
    start_stress_test()