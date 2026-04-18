<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键词">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="会员 / 备注 / 订单号"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="变动类型">
          <el-select v-model="searchInfo.changeType" clearable style="width: 160px">
            <el-option v-for="item in changeTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源类型">
          <el-select v-model="searchInfo.sourceType" clearable style="width: 180px">
            <el-option v-for="item in sourceTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
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
            {{ scope.row.member?.realName || scope.row.member?.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" min-width="120">
          <template #default="scope">
            {{ scope.row.member?.mobile || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="变动类型" min-width="110">
          <template #default="scope">
            <el-tag :type="tagTypeMap[scope.row.changeType] || ''">
              {{ changeTypeLabel(scope.row.changeType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="变动积分" min-width="100" prop="changePoints" />
        <el-table-column align="left" label="变动前" min-width="90" prop="beforePoints" />
        <el-table-column align="left" label="变动后" min-width="90" prop="afterPoints" />
        <el-table-column align="left" label="来源类型" min-width="140" prop="sourceType" />
        <el-table-column align="left" label="来源ID" min-width="100" prop="sourceId" />
        <el-table-column align="left" label="操作人ID" min-width="100" prop="operatorId" />
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
  import { getPointLogList } from '@/api/member'
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { useRoute } from 'vue-router'

  defineOptions({
    name: 'PointLogPage'
  })

  const route = useRoute()

  const changeTypeOptions = [
    { label: '获得积分', value: 'earn' },
    { label: '消耗积分', value: 'use' },
    { label: '手工增加', value: 'adjust_add' },
    { label: '手工扣减', value: 'adjust_sub' },
    { label: '退回积分', value: 'refund' }
  ]

  const sourceTypeOptions = [
    { label: '手工调整', value: 'manual' },
    { label: '兑换订单', value: 'exchange_order' },
    { label: '订单退回', value: 'exchange_order_void' }
  ]

  const tagTypeMap = {
    earn: 'success',
    use: '',
    adjust_add: 'success',
    adjust_sub: 'warning',
    refund: 'info'
  }

  const searchInfo = ref({
    memberId: route.query.memberId || '',
    keyword: '',
    changeType: '',
    sourceType: ''
  })

  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const changeTypeLabel = (value) => {
    const item = changeTypeOptions.find((option) => option.value === value)
    return item ? item.label : value || '-'
  }

  const getTableData = async () => {
    const res = await getPointLogList({
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
      changeType: '',
      sourceType: ''
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
