<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.mobile" clearable placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" clearable placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
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
        <el-button type="primary" icon="plus" @click="openCreateDialog">新增会员</el-button>
      </div>

      <el-table :data="tableData" row-key="id">
        <el-table-column align="left" label="创建时间" min-width="170">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" min-width="120" prop="mobile" />
        <el-table-column align="left" label="昵称" min-width="120" prop="nickname" />
        <el-table-column align="left" label="真实姓名" min-width="120" prop="realName" />
        <el-table-column align="left" label="会员等级" min-width="100" prop="memberLevel" />
        <el-table-column align="left" label="状态" min-width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'enabled' ? 'success' : 'danger'">
              {{ scope.row.status === 'enabled' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="OpenID" min-width="180" prop="openid" show-overflow-tooltip />
        <el-table-column align="left" label="操作" min-width="280" fixed="right">
          <template #default="scope">
            <el-button link type="primary" icon="edit" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button
              link
              :type="scope.row.status === 'enabled' ? 'warning' : 'success'"
              @click="handleStatusChange(scope.row)"
            >
              {{ scope.row.status === 'enabled' ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="primary" icon="tickets" @click="showAccount(scope.row)">积分账户</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogType === 'create' ? '新增会员' : '编辑会员'" width="720px">
      <el-form label-width="100px" :model="formData">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="OpenID">
              <el-input v-model="formData.openid" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="UnionID">
              <el-input v-model="formData.unionid" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input v-model="formData.mobile" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="昵称">
              <el-input v-model="formData.nickname" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="真实姓名">
              <el-input v-model="formData.realName" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="会员等级">
              <el-select v-model="formData.memberLevel" style="width: 100%">
                <el-option v-for="item in levelOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="formData.status" style="width: 100%">
                <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="头像地址">
              <el-input v-model="formData.avatarUrl" clearable placeholder="可直接填写图片 URL" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="accountDialogVisible" title="会员积分账户" width="520px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="会员">
          {{ currentAccount.member?.realName || currentAccount.member?.nickname || currentAccount.member?.mobile || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="手机号">
          {{ currentAccount.member?.mobile || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="可用积分">
          {{ currentAccount.availablePoints || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="冻结积分">
          {{ currentAccount.frozenPoints || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="累计获得">
          {{ currentAccount.totalEarnedPoints || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="累计使用">
          {{ currentAccount.totalUsedPoints || 0 }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    createMember,
    deleteMember,
    findMember,
    getMemberList,
    getMemberPointAccount,
    updateMember,
    updateMemberStatus
  } from '@/api/member'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { formatDate } from '@/utils/format'

  defineOptions({
    name: 'MemberListPage'
  })

  const defaultForm = () => ({
    id: 0,
    openid: '',
    unionid: '',
    mobile: '',
    nickname: '',
    avatarUrl: '',
    realName: '',
    memberLevel: 'standard',
    status: 'enabled'
  })

  const statusOptions = [
    { label: '启用', value: 'enabled' },
    { label: '禁用', value: 'disabled' }
  ]

  const levelOptions = [
    { label: '标准会员', value: 'standard' },
    { label: '轻奢会员', value: 'premium' },
    { label: '黑金会员', value: 'vip' }
  ]

  const searchInfo = ref({
    mobile: '',
    nickname: '',
    status: ''
  })
  const tableData = ref([])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  const dialogVisible = ref(false)
  const dialogType = ref('create')
  const formData = ref(defaultForm())

  const accountDialogVisible = ref(false)
  const currentAccount = ref({})

  const getTableData = async () => {
    const res = await getMemberList({
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

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const resetSearch = () => {
    searchInfo.value = {
      mobile: '',
      nickname: '',
      status: ''
    }
    page.value = 1
    getTableData()
  }

  const openCreateDialog = () => {
    dialogType.value = 'create'
    formData.value = defaultForm()
    dialogVisible.value = true
  }

  const openEditDialog = async (row) => {
    const res = await findMember({ id: row.id })
    if (res.code === 0) {
      dialogType.value = 'update'
      formData.value = {
        ...defaultForm(),
        ...res.data.member
      }
      dialogVisible.value = true
    }
  }

  const closeDialog = () => {
    dialogVisible.value = false
    formData.value = defaultForm()
  }

  const submitForm = async () => {
    const action = dialogType.value === 'create' ? createMember : updateMember
    const res = await action(formData.value)
    if (res.code === 0) {
      ElMessage.success(dialogType.value === 'create' ? '新增会员成功' : '更新会员成功')
      closeDialog()
      getTableData()
    }
  }

  const handleStatusChange = async (row) => {
    const nextStatus = row.status === 'enabled' ? 'disabled' : 'enabled'
    await ElMessageBox.confirm(`确定要${nextStatus === 'enabled' ? '启用' : '禁用'}该会员吗？`, '提示', {
      type: 'warning'
    })
    const res = await updateMemberStatus({
      id: row.id,
      status: nextStatus
    })
    if (res.code === 0) {
      ElMessage.success('会员状态更新成功')
      getTableData()
    }
  }

  const showAccount = async (row) => {
    const res = await getMemberPointAccount({ id: row.id })
    if (res.code === 0) {
      currentAccount.value = res.data
      accountDialogVisible.value = true
    }
  }

  const handleDelete = async (row) => {
    await ElMessageBox.confirm('删除会员会同时移除对应积分账户，确认继续吗？', '提示', {
      type: 'warning'
    })
    const res = await deleteMember({ id: row.id })
    if (res.code === 0) {
      ElMessage.success('删除会员成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }

  getTableData()
</script>

<style scoped>
  .gva-search-box {
    margin-bottom: 16px;
  }
</style>
