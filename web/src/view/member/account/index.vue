<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="会员检索">
          <el-input
            v-model="searchInfo.keyword"
            clearable
            placeholder="姓名 / 手机号 / 来源"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          <el-button icon="refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" row-key="id">
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
        <el-table-column align="left" label="可用积分" min-width="100" prop="balance" />
        <el-table-column align="left" label="冻结积分" min-width="100" prop="frozenPoints" />
        <el-table-column align="left" label="累计获得" min-width="100" prop="totalEarned" />
        <el-table-column align="left" label="累计消耗" min-width="100" prop="totalSpent" />
        <el-table-column align="left" label="更新时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="180" fixed="right">
          <template #default="scope">
            <el-button link type="success" icon="plus" @click="openAdjustDialog(scope.row, 'add')">加积分</el-button>
            <el-button link type="warning" icon="minus" @click="openAdjustDialog(scope.row, 'sub')">扣积分</el-button>
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

    <el-dialog v-model="adjustDialogVisible" :title="adjustType === 'add' ? '手工增加积分' : '手工扣减积分'" width="420px">
      <el-form label-width="90px" :model="adjustForm">
        <el-form-item label="当前会员">
          <div>{{ currentMemberLabel }}</div>
        </el-form-item>
        <el-form-item label="积分数量">
          <el-input-number v-model="adjustForm.points" :min="1" :step="10" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="adjustForm.remark" :rows="3" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeAdjustDialog">取消</el-button>
        <el-button type="primary" @click="submitAdjust">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { getPointAccountList, manualAddPoints, manualSubPoints } from '@/api/member'
  import { formatDate } from '@/utils/format'
  import { ElMessage } from 'element-plus'
  import { computed, ref } from 'vue'
  import { useRoute } from 'vue-router'

  defineOptions({
    name: 'PointAccountPage'
  })

  const route = useRoute()
  const searchInfo = ref({
    memberId: route.query.memberId || '',
    keyword: ''
  })

  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const adjustDialogVisible = ref(false)
  const adjustType = ref('add')
  const adjustForm = ref({
    memberId: 0,
    points: 10,
    remark: ''
  })
  const currentMember = ref(null)

  const currentMemberLabel = computed(() => {
    if (!currentMember.value) {
      return '-'
    }
    return currentMember.value.name || currentMember.value.phone || '-'
  })

  const getTableData = async () => {
    const res = await getPointAccountList({
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
      keyword: ''
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

  const openAdjustDialog = (row, type) => {
    currentMember.value = row.member
    adjustType.value = type
    adjustForm.value = {
      memberId: row.memberId,
      points: 10,
      remark: ''
    }
    adjustDialogVisible.value = true
  }

  const closeAdjustDialog = () => {
    adjustDialogVisible.value = false
    adjustForm.value = {
      memberId: 0,
      points: 10,
      remark: ''
    }
    currentMember.value = null
  }

  const submitAdjust = async () => {
    const action = adjustType.value === 'add' ? manualAddPoints : manualSubPoints
    const res = await action(adjustForm.value)
    if (res.code === 0) {
      ElMessage.success(adjustType.value === 'add' ? '积分增加成功' : '积分扣减成功')
      closeAdjustDialog()
      getTableData()
    }
  }

  getTableData()
</script>
