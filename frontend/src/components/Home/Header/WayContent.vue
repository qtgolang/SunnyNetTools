<script>
import {
  GetInterfaceOutRouterAdders,
  GetIPV4InterfaceAdders,
  SetInterfaceOutRouterAdders
} from "../../../../bindings/changeme/Service/appmain";
import {Tour_Add} from "../Tour";

export default {
  data() {
    return {
      WayContent: [],
    }
  },
  methods: {
    setOutRouter(ip) {
      SetInterfaceOutRouterAdders(ip).then(() => {
        for (let i = 0; i < this.WayContent.length; i++) {
          this.WayContent[i].outRouter = this.WayContent[i].ip === ip;
        }
      })
    },
    onPopoverShow() {
      this.WayContent = [{ip: '正在刷新', outRouter: false},]
      GetInterfaceOutRouterAdders().then(IPAdders => {
        GetIPV4InterfaceAdders().then(list => {
          const array = [];
          for (let i = 0; i < list.length; i++) {
            array.push({ip: list[i], outRouter: IPAdders === list[i]})
          }
          if (array.length === 0) {
            array.push({ip: '未获取到任何内网IP', outRouter: false})
          }
          this.WayContent = array
        })
      })
      /*
      this.WayContent = [
        {ip: '', outRouter: false},
        {ip: '暂未获取到', outRouter: false},
        {ip: '192.168.31.154', outRouter: true},
      ]
      */
    }
  },
  mounted() {
    Tour_Add(this.$refs.wg, 12, "网关查看", "点击这里可以查看所有网卡的IP\n\n方便您设置代理\n\n你也可以在设置数据出口IP、如果您有多出口网卡有用！")
  }
}
</script>

<template>
  <div ref="wg">
    <el-popover
        placement="top-start"
        :width="350"
        @show="onPopoverShow"
        trigger="hover"
        popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;"
    >
      <el-table :data="WayContent">
        <el-table-column width="150" label="当前内网IP">
          <template #default="scope">
            <div style="display: flex;justify-content: left;text-align: left">
              <el-text v-if="scope.row.ip === '正在刷新' || scope.row.ip === '未获取到任何内网IP'" type="primary"
                       size="small">
              </el-text>
              <el-text v-if="scope.row.ip !== '正在刷新' && scope.row.ip !== '未获取到任何内网IP'" type="primary"
                       size="small">
                {{ scope.row.ip }}
              </el-text>
            </div>
          </template>
        </el-table-column>
        <el-table-column width="150" label="操作">
          <template #default="scope">
            <div style="display: flex;justify-content: left;text-align: left">
              <el-text v-if="scope.row.ip === '正在刷新'" type="primary" size="small">
                正在刷新网关
              </el-text>
              <el-text v-if="scope.row.ip === '未获取到任何内网IP'" type="primary" size="small">
                未获取到任何内网IP
              </el-text>
              <el-button v-if="scope.row.outRouter&&scope.row.ip !== ''" type="primary" size="small"
                         @click="setOutRouter('')">
                取消指定为出口IP
              </el-button>
              <el-button v-if="!scope.row.outRouter&&scope.row.ip !== ''" type="primary" size="small"
                         @click="setOutRouter(scope.row.ip)">
                指定为出口IP
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <template #reference>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
          <path
              d="M24,0A24,24,0,1,0,48,24,24,24,0,0,0,24,0ZM35.32,35.8A40.15,40.15,0,0,0,37,25h9a21.91,21.91,0,0,1-5.85,13.94A24.94,24.94,0,0,0,35.32,35.8ZM8,39.1A21.92,21.92,0,0,1,2,25h9a40.05,40.05,0,0,0,1.71,10.94A25,25,0,0,0,8,39.1Zm4.64-26.81A40.21,40.21,0,0,0,11,23H2A21.91,21.91,0,0,1,7.85,9.08,25,25,0,0,0,12.66,12.28ZM25,15a24.92,24.92,0,0,0,8.51-1.85A38.76,38.76,0,0,1,35,23H25Zm0-2V2.1c3.2.61,6.05,4.1,7.88,9.12A22.9,22.9,0,0,1,25,13ZM23,2.1V13a22.91,22.91,0,0,1-7.88-1.74C17,6.19,19.8,2.71,23,2.1ZM23,15v8H13a38.75,38.75,0,0,1,1.48-9.87A24.93,24.93,0,0,0,23,15ZM13,25H23v8.2a24.9,24.9,0,0,0-8.44,1.89A38.63,38.63,0,0,1,13,25ZM23,35.23V45.9c-3.15-.6-6-4-7.8-8.9A22.89,22.89,0,0,1,23,35.23ZM25,45.9V35.22a22.93,22.93,0,0,1,7.85,1.66C31,41.85,28.18,45.3,25,45.9Zm0-12.7V25H35a38.7,38.7,0,0,1-1.51,10A24.94,24.94,0,0,0,25,33.2ZM37,23a40.21,40.21,0,0,0-1.64-10.72,24.94,24.94,0,0,0,4.8-3.21A21.91,21.91,0,0,1,46,23ZM38.71,7.66a23,23,0,0,1-4,2.71,21,21,0,0,0-4.5-7.48A22,22,0,0,1,38.71,7.66ZM13.3,10.36a23,23,0,0,1-4-2.71,22,22,0,0,1,8.52-4.76A21,21,0,0,0,13.3,10.36ZM9.47,40.5a23,23,0,0,1,3.92-2.65,20.82,20.82,0,0,0,4.42,7.25A22,22,0,0,1,9.47,40.5Zm25.2-2.79a23,23,0,0,1,4,2.65,22,22,0,0,1-8.5,4.75A21,21,0,0,0,34.67,37.71Z"
              fill="#0797E1"/>
        </svg>
      </template>
    </el-popover>
  </div>
</template>
