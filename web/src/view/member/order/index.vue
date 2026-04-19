<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="订单号 / 会员 / 手机号 / 商品 / 收货人"
            style="width: 300px"
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
        <el-button v-if="btnAuth.add" type="primary" icon="plus" @click="openCreateDialog">新建兑换订单</el-button>
      </div>

      <el-table :data="tableData" row-key="id">
        <el-table-column align="left" label="下单时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="订单号" min-width="180" prop="orderNo" />
        <el-table-column align="left" label="会员" min-width="120">
          <template #default="scope">
            {{ scope.row.member?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" min-width="120">
          <template #default="scope">
            {{ scope.row.member?.phone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="兑换商品" min-width="160" prop="productName" />
        <el-table-column align="left" label="数量" min-width="80" prop="quantity" />
        <el-table-column align="left" label="单件积分" min-width="100" prop="unitPoints" />
        <el-table-column align="left" label="总积分" min-width="100" prop="totalPoints" />
        <el-table-column align="left" label="收货人" min-width="120" prop="receiverName" />
        <el-table-column align="left" label="联系电话" min-width="120" prop="receiverPhone" />
        <el-table-column align="left" label="状态" min-width="100">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="备注" min-width="180" prop="remark" show-overflow-tooltip />
        <el-table-column v-if="showActionColumn" align="left" label="操作" min-width="240" fixed="right">
          <template #default="scope">
            <el-button v-if="btnAuth.info" link type="primary" icon="view" @click="openDetailDialog(scope.row)">详情</el-button>
            <el-button
              v-if="btnAuth.complete && scope.row.status === 1"
              link
              type="success"
              icon="select"
              @click="handleOrderAction(scope.row, 'complete')"
            >
              完成
            </el-button>
            <el-button
              v-if="btnAuth.cancel && scope.row.status === 1"
              link
              type="warning"
              icon="close-bold"
              @click="handleOrderAction(scope.row, 'cancel')"
            >
              取消
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

    <el-dialog v-model="createDialogVisible" title="新建兑换订单" width="620px">
      <el-form ref="createFormRef" label-width="100px" :model="createForm" :rules="rules">
        <el-form-item label="选择会员" prop="memberId">
          <el-select
            v-model="createForm.memberId"
            clearable
            filterable
            remote
            reserve-keyword
            placeholder="输入姓名 / 手机号检索会员"
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
        <el-form-item label="选择商品" prop="productId">
          <el-select
            v-model="createForm.productId"
            clearable
            filterable
            remote
            reserve-keyword
            placeholder="输入商品名称检索商品"
            style="width: 100%"
            :remote-method="loadProductOptions"
          >
            <el-option
              v-for="item in productOptions"
              :key="item.id"
              :label="`${item.name} · ${item.pointsPrice}积分 · 库存${item.stock}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="selectedProduct" label="商品信息">
          <div class="order-summary">
            <div>积分价格：{{ selectedProduct.pointsPrice }}</div>
            <div>当前库存：{{ selectedProduct.stock }}</div>
            <div>商品分类：{{ selectedProduct.category || '-' }}</div>
          </div>
        </el-form-item>
        <el-form-item label="兑换数量" prop="quantity">
          <el-input-number v-model="createForm.quantity" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="收货人" prop="receiverName">
          <el-input v-model="createForm.receiverName" clearable />
        </el-form-item>
        <el-form-item label="联系电话" prop="receiverPhone">
          <el-input v-model="createForm.receiverPhone" clearable />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
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
          {{ detailData.member?.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="手机号">
          {{ detailData.member?.phone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="兑换商品">
          {{ detailData.productName || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="商品ID">
          {{ detailData.productId || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="数量">
          {{ detailData.quantity || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="单件积分">
          {{ detailData.unitPoints || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="总积分">
          {{ detailData.totalPoints || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="收货人">
          {{ detailData.receiverName || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="联系电话">
          {{ detailData.receiverPhone || '-' }}
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
    cancelRedemptionOrder,
    completeRedemptionOrder,
    createRedemptionOrder,
    findRedemptionOrder,
    getMemberOptions,
    getPointProductOptions,
    getRedemptionOrderList
  } from '@/api/member'
  import { useBtnAuth } from '@/utils/btnAuth'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, nextTick, ref } from 'vue'

  defineOptions({
    name: 'RedemptionOrderPage'
  })

  const statusOptions = [
    { label: '待处理', value: 1 },
    { label: '已完成', value: 2 },
    { label: '已取消', value: 3 }
  ]

  const searchInfo = ref({
    keyword: '',
    status: ''
  })

  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const createFormRef = ref()
  const btnAuth = useBtnAuth()
  const createDialogVisible = ref(false)
  const detailDialogVisible = ref(false)
  const createForm = ref({
    memberId: '',
    productId: '',
    quantity: 1,
    receiverName: '',
    receiverPhone: '',
    remark: ''
  })
  const detailData = ref({})
  const memberOptions = ref([])
  const productOptions = ref([])

  const selectedProduct = computed(() => productOptions.value.find((item) => item.id === createForm.value.productId))
  const showActionColumn = computed(() => Boolean(btnAuth.info || btnAuth.complete || btnAuth.cancel))
  const phonePattern = /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/
  const trimmedRequired = (message) => ({
    validator: (_, value, callback) => {
      if (typeof value === 'string' ? value.trim() : value) {
        callback()
        return
      }
      callback(new Error(message))
    },
    trigger: 'blur'
  })
  const quantityValidator = (_, value, callback) => {
    if (typeof value !== 'number' || value <= 0) {
      callback(new Error('兑换数量必须大于 0'))
      return
    }
    if (selectedProduct.value && value > selectedProduct.value.stock) {
      callback(new Error('兑换数量不能超过当前库存'))
      return
    }
    callback()
  }
  const rules = {
    memberId: [{ required: true, message: '请选择会员', trigger: 'change' }],
    productId: [{ required: true, message: '请选择商品', trigger: 'change' }],
    quantity: [{ validator: quantityValidator, trigger: 'change' }],
    receiverName: [trimmedRequired('请输入收货人')],
    receiverPhone: [
      trimmedRequired('请输入联系电话'),
      { pattern: phonePattern, message: '请输入合法手机号', trigger: 'blur' }
    ],
    remark: [{ max: 200, message: '备注不能超过 200 个字符', trigger: 'blur' }]
  }

  const statusLabel = (value) => {
    const item = statusOptions.find((option) => option.value === value)
    return item ? item.label : value || '-'
  }

  const statusTagType = (value) => {
    switch (value) {
      case 1:
        return 'warning'
      case 2:
        return 'success'
      case 3:
        return 'info'
      default:
        return ''
    }
  }

  const formatMemberLabel = (item) => {
    const label = item.name || item.phone || '-'
    return `${label} · ${item.phone || '未填写手机号'}`
  }

  const getTableData = async () => {
    const res = await getRedemptionOrderList({
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

  const loadProductOptions = async (keyword = '') => {
    const res = await getPointProductOptions({ keyword })
    if (res.code === 0) {
      productOptions.value = res.data.list || []
    }
  }

  const openCreateDialog = async () => {
    createForm.value = {
      memberId: '',
      productId: '',
      quantity: 1,
      receiverName: '',
      receiverPhone: '',
      remark: ''
    }
    await Promise.all([loadMemberOptions(), loadProductOptions()])
    createDialogVisible.value = true
    await nextTick()
    createFormRef.value?.clearValidate()
  }

  const closeCreateDialog = () => {
    createDialogVisible.value = false
    createForm.value = {
      memberId: '',
      productId: '',
      quantity: 1,
      receiverName: '',
      receiverPhone: '',
      remark: ''
    }
    createFormRef.value?.clearValidate()
  }

  const submitCreate = async () => {
    if (!createFormRef.value) {
      return
    }
    try {
      await createFormRef.value.validate()
    } catch {
      return
    }
    const res = await createRedemptionOrder(createForm.value)
    if (res.code === 0) {
      ElMessage.success('兑换订单创建成功')
      closeCreateDialog()
      getTableData()
    }
  }

  const openDetailDialog = async (row) => {
    const res = await findRedemptionOrder({ id: row.id })
    if (res.code === 0) {
      detailData.value = res.data
      detailDialogVisible.value = true
    }
  }

  const handleOrderAction = async (row, action) => {
    const actionMap = {
      complete: {
        message: '确认将该订单标记为已完成吗？',
        request: completeRedemptionOrder,
        success: '订单已完成'
      },
      cancel: {
        message: '确认取消该订单并退回积分吗？',
        request: cancelRedemptionOrder,
        success: '订单已取消'
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
