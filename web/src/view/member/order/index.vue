<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="订单号 / 核销码 / 会员 / 商品"
            style="width: 280px"
          />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchInfo.status" clearable style="width: 180px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          <el-button icon="refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreateDialog">新建兑换订单</el-button>
      </div>

      <el-table :data="tableData" row-key="id">
        <el-table-column align="left" label="下单时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="订单号" min-width="180" prop="orderNo" />
        <el-table-column align="left" label="会员" min-width="150">
          <template #default="scope">
            {{ scope.row.member?.realName || scope.row.member?.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" min-width="120">
          <template #default="scope">
            {{ scope.row.member?.mobile || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="兑换商品" min-width="160">
          <template #default="scope">
            {{ scope.row.goods?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="消耗积分" min-width="100" prop="pointsCost" />
        <el-table-column align="left" label="状态" min-width="100">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="核销码" min-width="100" prop="verifyCode" />
        <el-table-column align="left" label="核销时间" min-width="170">
          <template #default="scope">
            {{ scope.row.verifiedAt ? formatDate(scope.row.verifiedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="280" fixed="right">
          <template #default="scope">
            <el-button link type="primary" icon="view" @click="openDetailDialog(scope.row)">详情</el-button>
            <el-button
              v-if="scope.row.status === 'pending'"
              link
              type="success"
              icon="select"
              @click="handleOrderAction(scope.row, 'verify')"
            >
              核销
            </el-button>
            <el-button
              v-if="scope.row.status === 'pending'"
              link
              type="warning"
              icon="close-bold"
              @click="handleOrderAction(scope.row, 'cancel')"
            >
              取消
            </el-button>
            <el-button
              v-if="scope.row.status === 'completed'"
              link
              type="danger"
              icon="refresh-left"
              @click="handleOrderAction(scope.row, 'refund')"
            >
              退款
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="createDialogVisible" title="新建兑换订单" width="560px">
      <el-form label-width="100px" :model="createForm">
        <el-form-item label="选择会员">
          <el-select
            v-model="createForm.memberId"
            clearable
            filterable
            remote
            reserve-keyword
            placeholder="输入手机号 / 昵称检索会员"
            style="width: 100%"
            :remote-method="loadMemberOptions"
          >
            <el-option
              v-for="item in memberOptions"
              :key="item.id"
              :label="formatMemberLabel(item)"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择商品">
          <el-select
            v-model="createForm.goodsId"
            clearable
            filterable
            remote
            reserve-keyword
            placeholder="输入商品名称检索商品"
            style="width: 100%"
            :remote-method="loadGoodsOptions"
          >
            <el-option
              v-for="item in goodsOptions"
              :key="item.id"
              :label="`${item.name} · ${item.pointsPrice}积分 · 库存${item.stock}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="selectedGoods" label="商品信息">
          <div class="order-summary">
            <div>积分价格：{{ selectedGoods.pointsPrice }}</div>
            <div>当前库存：{{ selectedGoods.stock }}</div>
            <div>每人限兑：{{ selectedGoods.limitPerMember || '不限' }}</div>
          </div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" :rows="3" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeCreateDialog">取消</el-button>
        <el-button type="primary" @click="submitCreate">确认兑换</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="订单详情" width="720px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">
          {{ detailData.orderNo || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="订单状态">
          {{ statusLabel(detailData.status) }}
        </el-descriptions-item>
        <el-descriptions-item label="会员">
          {{ detailData.member?.realName || detailData.member?.nickname || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="手机号">
          {{ detailData.member?.mobile || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="兑换商品">
          {{ detailData.goods?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="消耗积分">
          {{ detailData.pointsCost || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="核销码">
          {{ detailData.verifyCode || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="核销时间">
          {{ detailData.verifiedAt ? formatDate(detailData.verifiedAt) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ detailData.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    cancelExchangeOrder,
    createExchangeOrder,
    findExchangeOrder,
    getExchangeOrderList,
    getMemberOptions,
    getPointGoodsOptions,
    refundExchangeOrder,
    verifyExchangeOrder
  } from '@/api/member'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, ref } from 'vue'

  defineOptions({
    name: 'ExchangeOrderPage'
  })

  const statusOptions = [
    { label: '待核销', value: 'pending' },
    { label: '已完成', value: 'completed' },
    { label: '已取消', value: 'cancelled' },
    { label: '已退款', value: 'refunded' }
  ]

  const searchInfo = ref({
    keyword: '',
    status: ''
  })

  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const createDialogVisible = ref(false)
  const detailDialogVisible = ref(false)
  const createForm = ref({
    memberId: '',
    goodsId: '',
    remark: ''
  })
  const detailData = ref({})
  const memberOptions = ref([])
  const goodsOptions = ref([])

  const selectedGoods = computed(() => goodsOptions.value.find((item) => item.id === createForm.value.goodsId))

  const statusLabel = (value) => {
    const item = statusOptions.find((option) => option.value === value)
    return item ? item.label : value || '-'
  }

  const statusTagType = (value) => {
    switch (value) {
      case 'pending':
        return 'warning'
      case 'completed':
        return 'success'
      case 'cancelled':
        return 'info'
      case 'refunded':
        return 'danger'
      default:
        return ''
    }
  }

  const formatMemberLabel = (item) => {
    const label = item.realName || item.nickname || item.mobile
    return `${label} · ${item.mobile || '未绑定手机号'}`
  }

  const getTableData = async () => {
    const res = await getExchangeOrderList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.pageSize
    }
  }

  const resetSearch = () => {
    searchInfo.value = {
      keyword: '',
      status: ''
    }
    page.value = 1
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const loadMemberOptions = async (keyword = '') => {
    const res = await getMemberOptions({ keyword })
    if (res.code === 0) {
      memberOptions.value = res.data.list || []
    }
  }

  const loadGoodsOptions = async (keyword = '') => {
    const res = await getPointGoodsOptions({ keyword })
    if (res.code === 0) {
      goodsOptions.value = res.data.list || []
    }
  }

  const openCreateDialog = async () => {
    createForm.value = {
      memberId: '',
      goodsId: '',
      remark: ''
    }
    await Promise.all([loadMemberOptions(), loadGoodsOptions()])
    createDialogVisible.value = true
  }

  const closeCreateDialog = () => {
    createDialogVisible.value = false
    createForm.value = {
      memberId: '',
      goodsId: '',
      remark: ''
    }
  }

  const submitCreate = async () => {
    const res = await createExchangeOrder(createForm.value)
    if (res.code === 0) {
      ElMessage.success('兑换订单创建成功')
      closeCreateDialog()
      getTableData()
    }
  }

  const openDetailDialog = async (row) => {
    const res = await findExchangeOrder({ id: row.id })
    if (res.code === 0) {
      detailData.value = res.data
      detailDialogVisible.value = true
    }
  }

  const handleOrderAction = async (row, action) => {
    const actionMap = {
      verify: {
        message: '确认核销该订单吗？',
        request: verifyExchangeOrder,
        success: '订单核销成功'
      },
      cancel: {
        message: '确认取消该订单并退回积分吗？',
        request: cancelExchangeOrder,
        success: '订单取消成功'
      },
      refund: {
        message: '确认退回积分并退款该订单吗？',
        request: refundExchangeOrder,
        success: '订单退款成功'
      }
    }
    const config = actionMap[action]
    await ElMessageBox.confirm(config.message, '提示', {
      type: 'warning'
    })
    const res = await config.request({ id: row.id, remark: '' })
    if (res.code === 0) {
      ElMessage.success(config.success)
      getTableData()
    }
  }

  getTableData()
</script>

<style scoped>
  .order-summary {
    line-height: 28px;
  }
</style>
