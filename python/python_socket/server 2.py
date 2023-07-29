import socket
import time
import threading

NUMBEROFCONNECTIONS = 0


def main():
    for j in range(1):
        for i in range(2000):
            threading.Thread(
                target=server_program,
                args=(socket.gethostbyaddr(f"127.0.0.1")[0], 8000 + i),
            ).start()


def server_program(host, port):
    # get the hostname
    # socket.gethostname()
    # host = socket.gethostbyaddr("127.0.0.1")[0]
    # print(host)
    # port = 8000  # initiate port no above 1024
    print("Listening to:", host, port)

    server_socket = socket.socket()  # get instance
    # look closely. The bind() function takes tuple as argument
    server_socket.bind((host, port))  # bind host address and port together

    # configure how many client the server can listen simultaneously
    server_socket.listen(2)
    print("Listening to:", server_socket.getsockname())
    conn, address = server_socket.accept()  # accept new connection

    print("Connection from: " + str(address))
    while True:
        # receive data stream. it won't accept data packet greater than 1024 bytes
        data = conn.recv(1024).decode()
        print("from connected user: " + str(data))

        # current time
        current_time = time.ctime(time.time()) + "\r\n"
        # convert to bytes
        current_time = current_time.encode()
        conn.send(current_time)
        time.sleep(60)
        # send data to the client

    conn.close()  # close the connection


if __name__ == "__main__":
    main()
