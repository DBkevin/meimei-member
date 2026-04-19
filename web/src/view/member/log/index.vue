<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="会员 / 备注 / 操作人"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="流水类型">
          <el-select v-model="searchInfo.type" clearable style="width: 160px">
            <el-option v-for="item in typeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源类型">
          <el-select v-model="searchInfo.refType" clearable style="width: 200px">
            <el-option v-for="item in refTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          <el-button icon="refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id">
        <el-table-column align="left" label="时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="会员" min-width="150">
          <template #default="scope">
            {{ scope.row.member?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" min-width="120">
          <template #default="scope">
            {{ scope.row.member?.phone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="流水类型" min-width="110">
          <template #default="scope">
            <el-tag :type="tagTypeMap[scope.row.type] || ''">
              {{ typeLabel(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="积分数量" min-width="100" prop="points" />
        <el-table-column align="left" label="变动前" min-width="90" prop="beforeBalance" />
        <el-table-column align="left" label="变动后" min-width="90" prop="afterBalance" />
        <el-table-column align="left" label="来源类型" min-width="170">
          <template #default="scope">
            {{ refTypeLabel(scope.row.refType) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="来源ID" min-width="100" prop="refId" />
        <el-table-column align="left" label="操作人" min-width="100" prop="operator" />
        <el-table-column align="left" label="备注" min-width="220" prop="remark" show-overflow-tooltip />
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
  </div>
</template>

<script setup>
  import { getPointTransactionList } from '@/api/member'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { useRoute } from 'vue-router'

  defineOptions({
    name: 'PointTransactionPage'
  })

  const route = useRoute()

  const typeOptions = [
    { label: '获得积分', value: 'earn' },
    { label: '消耗积分', value: 'spend' },
    { label: '手工调整', value: 'adjust' },
    { label: '退回积分', value: 'refund' }
  ]

  const refTypeOptions = [
    { label: '手工增加', value: 'manual_adjust_add' },
    { label: '手工扣减', value: 'manual_adjust_sub' },
    { label: '兑换订单', value: 'redemption_order' },
    { label: '订单取消退回', value: 'redemption_order_cancel' }
  ]

  const tagTypeMap = {
    earn: 'success',
    spend: 'warning',
    adjust: 'info',
    refund: 'primary'
  }

  const searchInfo = ref({
    memberId: route.query.memberId || '',
    keyword: '',
    type: '',
    refType: ''
  })

  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const typeLabel = (value) => {
    const item = typeOptions.find((option) => option.value === value)
    return item ? item.label : value || '-'
  }

  const refTypeLabel = (value) => {
    const item = refTypeOptions.find((option) => option.value === value)
    return item ? item.label : value || '-'
  }

  const getTableData = async () => {
    const res = await getPointTransactionList({
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
      memberId: route.query.memberId || '',
      keyword: '',
      type: '',
      refType: ''
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

  getTableData()
</script>
