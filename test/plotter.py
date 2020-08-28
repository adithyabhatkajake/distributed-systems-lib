import matplotlib.pyplot as plt
from data import *

def plot_different_block_num_commands():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Number of Commands")
    y = times
    x1 = b5_synchs_cmds
    x2 = b10_synchs_cmds
    x3 = b20_synchs_cmds
    X1 = b5_apollo_cmds
    X2 = b10_apollo_cmds
    X3 = b20_apollo_cmds
    plt.plot(y,x1, label="SyncHS Block size 5")
    plt.plot(y,x2, label="SyncHS Block size 10")
    plt.plot(y,x3, label="SyncHS Block size 20")
    plt.plot(y,X1, label="Apollo Block size 5")
    plt.plot(y,X2, label="Apollo Block size 10")
    plt.plot(y,X3, label="Apollo Block size 20")
    plt.legend()
    plt.savefig("different_block_num_commands.pdf")
    plt.clf()

def plot_different_block_num_latencies():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Latency (in ms)")
    y = times
    x1 = [b5_synchs_total_times[i]/b5_synchs_cmds[i] for i in range(5)]
    x2 = [b10_synchs_total_times[i]/b10_synchs_cmds[i] for i in range(5)]
    x3 = [b20_synchs_total_times[i]/b20_synchs_cmds[i] for i in range(5)]
    X1 = [b5_synchs_total_times[i]/b5_synchs_cmds[i] for i in range(5)]
    X2 = [b10_synchs_total_times[i]/b10_synchs_cmds[i] for i in range(5)]
    X3 = [b20_synchs_total_times[i]/b20_synchs_cmds[i] for i in range(5)]
    plt.plot(y,x1, label="SyncHS Block size 5")
    plt.plot(y,x2, label="SyncHS Block size 10")
    plt.plot(y,x3, label="SyncHS Block size 20")
    plt.plot(y,X1, label="Apollo Block size 5")
    plt.plot(y,X2, label="Apollo Block size 10")
    plt.plot(y,X3, label="Apollo Block size 20")
    plt.legend()
    plt.savefig("different_block_num_latencies.pdf")
    plt.clf()

def plot_different_delta_num_commands():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Number of Commands")
    y = times
    x1 = d2_f1_synchs_cmds
    x2 = d5_f1_synchs_cmds
    X1 = d2_f1_apollo_cmds
    X2 = d5_f1_apollo_cmds
    plt.plot(y,x1, label="SyncHS Delta 2")
    plt.plot(y,x2, label="SyncHS Delta 5")
    plt.plot(y,X1, label="Apollo Delta 2")
    plt.plot(y,X2, label="Apollo Delta 5")
    plt.legend()
    plt.savefig("different_delta_num_commands.pdf")
    plt.clf()

def plot_different_delta_num_latencies():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Latency (in ms)")
    y = times
    x1 = [d2_f1_synchs_total_times[i]/d2_f1_synchs_cmds[i] for i in range(5)]
    x2 = [d5_f1_synchs_total_times[i]/d5_f1_synchs_cmds[i] for i in range(5)]
    X1 = [d2_f1_apollo_total_times[i]/d2_f1_apollo_cmds[i] for i in range(5)]
    X2 = [d5_f1_apollo_total_times[i]/d5_f1_apollo_cmds[i] for i in range(5)]
    plt.plot(y,x1, label="SyncHS Delta 2")
    plt.plot(y,x2, label="SyncHS Delta 5")
    plt.plot(y,X1, label="Apollo Delta 2")
    plt.plot(y,X2, label="Apollo Delta 5")
    plt.legend()
    plt.savefig("different_delta_num_latencies.pdf")
    plt.clf()


def plot_different_f_num_commands():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Number of Commands")
    y = times
    x1 = d2_f1_synchs_cmds
    x2 = d2_f2_synchs_cmds
    X1 = d2_f1_apollo_cmds
    X2 = d2_f2_apollo_cmds
    plt.plot(y,x1, label="SyncHS f 1")
    plt.plot(y,x2, label="SyncHS f 2")
    plt.plot(y,X1, label="Apollo f 1")
    plt.plot(y,X2, label="Apollo f 2")
    plt.legend()
    plt.savefig("different_f_num_commands.pdf")
    plt.clf()


def plot_different_f_num_latencies():
    plt.xlabel("Time(in seconds)")
    plt.ylabel("Latency (in ms)")
    y = times
    x1 = [d2_f1_synchs_total_times[i]/d2_f1_synchs_cmds[i] for i in range(5)]
    x2 = [d2_f2_synchs_total_times[i]/d2_f2_synchs_cmds[i] for i in range(5)]
    X1 = [d2_f1_apollo_total_times[i]/d2_f1_apollo_cmds[i] for i in range(5)]
    X2 = [d2_f2_apollo_total_times[i]/d2_f2_apollo_cmds[i] for i in range(5)]
    plt.plot(y,x1, label="SyncHS f 1")
    plt.plot(y,x2, label="SyncHS f 2")
    plt.plot(y,X1, label="Apollo f 1")
    plt.plot(y,X2, label="Apollo f 2")
    plt.legend()
    plt.savefig("different_f_num_latencies.pdf")
    plt.clf()

if __name__ == "__main__":
    plot_different_block_num_commands()
    plot_different_block_num_latencies()
    plot_different_delta_num_commands()
    plot_different_delta_num_latencies()
    plot_different_f_num_commands()
    plot_different_f_num_latencies()
