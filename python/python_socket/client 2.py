import socket

import threading


class Host:
    def __init__(self):
        # 127.0.0.1
        self.first_chunk = 127
        self.second_chunk = 0
        self.third_chunk = 0
        self.fourth_chunk = 0

    def upgrade_host(self):
        if self.fourth_chunk < 255:
            self.fourth_chunk += 1
        elif self.third_chunk < 255:
            self.third_chunk += 1
            self.fourth_chunk = 0
        elif self.second_chunk < 255:
            self.second_chunk += 1
            self.third_chunk = 0
            self.fourth_chunk = 0
        else:
            print("No more hosts available")
            return False
        return True

    def __str__(self):
        return f"{self.first_chunk}.{self.second_chunk}.{self.third_chunk}.{self.fourth_chunk}"


def main():
    NUMBEROFCONNECTIONTOBEMADE = input("Number Of Connections to made: ")

    host = Host()

    for i in range(int(NUMBEROFCONNECTIONTOBEMADE) // 2000):
        if host.upgrade_host():
            for i in range(2000):
                print(host)
                threading.Thread(
                    target=client_program,
                    args=(socket.gethostbyaddr(str(host))[0], 8000 + i),
                ).start()
                # client_program(socket.gethostbyaddr("127.0.0.1")[0], 8000 + i)
        else:
            break  # No more hosts available

    while True:
        continue


def client_program(host, port):
    # host = socket.gethostbyaddr("127.0.0.1")[0]  # as both code is running on same pc
    # port = 8000  # socket server port number

    client_socket = socket.socket()  # instantiate
    client_socket.connect((host, port))  # connect to the server
    print(host)

    message = "Hello world"  # take input

    while message.lower().strip() != "bye":
        client_socket.send(message.encode())  # send message
        data = client_socket.recv(1024).decode()  # receive response

        print("Received from server: " + data)  # show in terminal

        # message = input(" -> ")  # again take input

    client_socket.close()  # close the connection


if __name__ == "__main__":
    main()
