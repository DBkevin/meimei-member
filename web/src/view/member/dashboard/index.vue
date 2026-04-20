<template>
  <div class="dashboard-container">
    <!-- 顶部统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409EFF">
              <i class="el-icon-user"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ summary.totalMembers }}</div>
              <div class="stat-label">会员总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67C23A">
              <i class="el-icon-coin"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ summary.totalPointsBalance }}</div>
              <div class="stat-label">当前总积分</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #E6A23C">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ summary.todayNewMembers }}</div>
              <div class="stat-label">今日新增会员</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #F56C6C">
              <i class="el-icon-shopping-cart-2"></i>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ summary.pendingOrders }}</div>
              <div class="stat-label">待处理订单</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二部分：订单和商品状态 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">
            <span>订单状态统计</span>
          </div>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="待处理">{{ summary.pendingOrders }}</el-descriptions-item>
            <el-descriptions-item label="已完成">{{ summary.completedOrders }}</el-descriptions-item>
            <el-descriptions-item label="已取消">{{ summary.cancelledOrders }}</el-descriptions-item>
            <el-descriptions-item label="今日新增">{{ summary.todayNewOrders }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">
            <span>商品状态统计</span>
          </div>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="上架中">{{ summary.onSaleProducts }}</el-descriptions-item>
            <el-descriptions-item label="已下架">{{ summary.offSaleProducts }}</el-descriptions-item>
            <el-descriptions-item label="库存不足">
              <el-tag type="danger">{{ summary.outOfStockProducts }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="商品总数">{{ summary.totalProducts }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第三部分：最近流水和订单 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">
            <span>最近积分流水</span>
          </div>
          <el-table :data="recentTransactions" size="small" max-height="300">
            <el-table-column prop="memberId" label="会员ID" width="80"></el-table-column>
            <el-table-column prop="type" label="类型" width="100">
              <template slot-scope="scope">
                <el-tag :type="getTransactionTagType(scope.row.type)">
                  {{ getTransactionTypeName(scope.row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="积分" width="80"></el-table-column>
            <el-table-column prop="afterBalance" label="余额" width="80"></el-table-column>
            <el-table-column prop="createdAt" label="时间" width="160"></el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">
            <span>最近兑换订单</span>
          </div>
          <el-table :data="recentOrders" size="small" max-height="300">
            <el-table-column prop="orderNo" label="订单号" width="160"></el-table-column>
            <el-table-column prop="memberId" label="会员ID" width="80"></el-table-column>
            <el-table-column prop="totalPoints" label="积分" width="80"></el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template slot-scope="scope">
                <el-tag :type="getOrderTagType(scope.row.status)">
                  {{ getOrderStatusName(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 库存不足商品提醒 -->
    <el-row :gutter="20" style="margin-top: 20px" v-if="lowStockProducts.length > 0">
      <el-col :span="24">
        <el-card shadow="hover">
          <div slot="header">
            <span style="color: #F56C6C">⚠️ 库存不足商品提醒</span>
          </div>
          <el-table :data="lowStockProducts" size="small">
            <el-table-column prop="id" label="ID" width="60"></el-table-column>
            <el-table-column prop="name" label="商品名称"></el-table-column>
            <el-table-column prop="category" label="分类" width="100"></el-table-column>
            <el-table-column prop="pointsPrice" label="积分价格" width="100"></el-table-column>
            <el-table-column prop="stock" label="库存" width="80">
              <template slot-scope="scope">
                <el-tag type="danger">{{ scope.row.stock }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template slot-scope="scope">
                <el-tag>{{ scope.row.status === 1 ? '上架' : '下架' }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getDashboardSummary, getRecentTransactions, getRecentOrders, getLowStockProducts } from '@/api/member'

export default {
  name: 'Dashboard',
  data() {
    return {
      summary: {
        totalMembers: 0,
        enabledMembers: 0,
        disabledMembers: 0,
        todayNewMembers: 0,
        totalPointsBalance: 0,
        totalPointsIssued: 0,
        totalPointsConsumed: 0,
        todayPointsIssued: 0,
        todayPointsConsumed: 0,
        totalProducts: 0,
        onSaleProducts: 0,
        offSaleProducts: 0,
        outOfStockProducts: 0,
        totalOrders: 0,
        pendingOrders: 0,
        completedOrders: 0,
        cancelledOrders: 0,
        todayNewOrders: 0
      },
      recentTransactions: [],
      recentOrders: [],
      lowStockProducts: []
    }
  },
  created() {
    this.loadData()
  },
  methods: {
    async loadData() {
      try {
        const [summaryRes, transactionsRes, ordersRes, productsRes] = await Promise.all([
          getDashboardSummary(),
          getRecentTransactions({ limit: 10 }),
          getRecentOrders({ limit: 10 }),
          getLowStockProducts({ limit: 10 })
        ])

        if (summaryRes.data.code === 0) {
          this.summary = summaryRes.data.data
        }
        if (transactionsRes.data.code === 0) {
          this.recentTransactions = transactionsRes.data.data || []
        }
        if (ordersRes.data.code === 0) {
          this.recentOrders = ordersRes.data.data || []
        }
        if (productsRes.data.code === 0) {
          this.lowStockProducts = productsRes.data.data || []
        }
      } catch (error) {
        this.$message.error('获取数据失败：' + error.message)
      }
    },
    getTransactionTypeName(type) {
      const map = {
        'earn': '增加积分',
        'spend': '消耗积分',
        'adjust': '手动调整',
        'refund': '退款返还'
      }
      return map[type] || type
    },
    getTransactionTagType(type) {
      const map = {
        'earn': 'success',
        'spend': 'warning',
        'adjust': 'info',
        'refund': 'danger'
      }
      return map[type] || 'info'
    },
    getOrderStatusName(status) {
      const map = {
        1: '待处理',
        2: '已完成',
        3: '已取消'
      }
      return map[status] || status
    },
    getOrderTagType(status) {
      const map = {
        1: 'warning',
        2: 'success',
        3: 'info'
      }
      return map[status] || 'info'
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}
.stat-cards {
  margin-bottom: 20px;
}
.stat-card {
  display: flex;
  align-items: center;
}
.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #fff;
  margin-right: 15px;
}
.stat-content {
  flex: 1;
}
.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}
.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}
</style>