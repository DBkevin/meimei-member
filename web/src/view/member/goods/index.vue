<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="商品关键词">
          <el-input v-model="searchInfo.keyword" clearable placeholder="请输入商品名称" style="width: 240px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable style="width: 160px">
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
        <el-button type="primary" icon="plus" @click="openCreateDialog">新增商品</el-button>
      </div>

      <el-table :data="tableData" row-key="id">
        <el-table-column align="left" label="封面" width="96">
          <template #default="scope">
            <el-image
              v-if="scope.row.coverImage"
              :src="scope.row.coverImage"
              fit="cover"
              style="width: 56px; height: 56px; border-radius: 12px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="商品名称" min-width="160" prop="name" />
        <el-table-column align="left" label="积分价格" min-width="100" prop="pointsPrice" />
        <el-table-column align="left" label="库存" min-width="90" prop="stock" />
        <el-table-column align="left" label="每人限兑" min-width="100" prop="limitPerMember" />
        <el-table-column align="left" label="状态" min-width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'on_sale' ? 'success' : 'info'">
              {{ scope.row.status === 'on_sale' ? '上架中' : '已下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" min-width="80" prop="sort" />
        <el-table-column align="left" label="更新时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="260" fixed="right">
          <template #default="scope">
            <el-button link type="primary" icon="edit" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button link type="primary" icon="box" @click="openStockDialog(scope.row)">库存</el-button>
            <el-button
              link
              :type="scope.row.status === 'on_sale' ? 'warning' : 'success'"
              @click="toggleStatus(scope.row)"
            >
              {{ scope.row.status === 'on_sale' ? '下架' : '上架' }}
            </el-button>
            <el-button link type="danger" icon="delete" @click="handleDelete(scope.row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogType === 'create' ? '新增积分商品' : '编辑积分商品'" width="760px">
      <el-form label-width="100px" :model="formData">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="商品名称">
              <el-input v-model="formData.name" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="积分价格">
              <el-input-number v-model="formData.pointsPrice" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="库存">
              <el-input-number v-model="formData.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="每人限兑">
              <el-input-number v-model="formData.limitPerMember" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="formData.status" style="width: 100%">
                <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="封面地址">
              <el-input v-model="formData.coverImage" clearable placeholder="可直接填写图片 URL" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="商品描述">
              <el-input v-model="formData.description" :rows="5" type="textarea" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="stockDialogVisible" title="库存管理" width="420px">
      <el-form label-width="90px" :model="stockForm">
        <el-form-item label="商品名称">
          <div>{{ stockGoodsName }}</div>
        </el-form-item>
        <el-form-item label="库存数量">
          <el-input-number v-model="stockForm.stock" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeStockDialog">取消</el-button>
        <el-button type="primary" @click="submitStock">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    createPointGoods,
    deletePointGoods,
    findPointGoods,
    getPointGoodsList,
    updatePointGoods,
    updatePointGoodsStatus,
    updatePointGoodsStock
  } from '@/api/member'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { computed, ref } from 'vue'

  defineOptions({
    name: 'PointGoodsPage'
  })

  const defaultForm = () => ({
    id: 0,
    name: '',
    coverImage: '',
    description: '',
    pointsPrice: 100,
    stock: 0,
    limitPerMember: 0,
    status: 'on_sale',
    sort: 0
  })

  const statusOptions = [
    { label: '上架中', value: 'on_sale' },
    { label: '已下架', value: 'off_sale' }
  ]

  const searchInfo = ref({
    keyword: '',
    status: ''
  })
  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const dialogVisible = ref(false)
  const dialogType = ref('create')
  const formData = ref(defaultForm())

  const stockDialogVisible = ref(false)
  const stockForm = ref({
    id: 0,
    stock: 0
  })
  const stockGoodsName = computed(() => {
    const current = tableData.value.find((item) => item.id === stockForm.value.id)
    return current?.name || '-'
  })

  const getTableData = async () => {
    const res = await getPointGoodsList({
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

  const openCreateDialog = () => {
    dialogType.value = 'create'
    formData.value = defaultForm()
    dialogVisible.value = true
  }

  const openEditDialog = async (row) => {
    const res = await findPointGoods({ id: row.id })
    if (res.code === 0) {
      dialogType.value = 'update'
      formData.value = {
        ...defaultForm(),
        ...res.data
      }
      dialogVisible.value = true
    }
  }

  const closeDialog = () => {
    dialogVisible.value = false
    formData.value = defaultForm()
  }

  const submitForm = async () => {
    const action = dialogType.value === 'create' ? createPointGoods : updatePointGoods
    const res = await action(formData.value)
    if (res.code === 0) {
      ElMessage.success(dialogType.value === 'create' ? '新增商品成功' : '更新商品成功')
      closeDialog()
      getTableData()
    }
  }

  const toggleStatus = async (row) => {
    const nextStatus = row.status === 'on_sale' ? 'off_sale' : 'on_sale'
    const res = await updatePointGoodsStatus({
      id: row.id,
      status: nextStatus
    })
    if (res.code === 0) {
      ElMessage.success(nextStatus === 'on_sale' ? '商品已上架' : '商品已下架')
      getTableData()
    }
  }

  const openStockDialog = (row) => {
    stockForm.value = {
      id: row.id,
      stock: row.stock
    }
    stockDialogVisible.value = true
  }

  const closeStockDialog = () => {
    stockDialogVisible.value = false
    stockForm.value = {
      id: 0,
      stock: 0
    }
  }

  const submitStock = async () => {
    const res = await updatePointGoodsStock(stockForm.value)
    if (res.code === 0) {
      ElMessage.success('库存更新成功')
      closeStockDialog()
      getTableData()
    }
  }

  const handleDelete = async (row) => {
    await ElMessageBox.confirm('删除商品后不可恢复，确认继续吗？', '提示', {
      type: 'warning'
    })
    const res = await deletePointGoods({ id: row.id })
    if (res.code === 0) {
      ElMessage.success('删除商品成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }

  getTableData()
</script>
