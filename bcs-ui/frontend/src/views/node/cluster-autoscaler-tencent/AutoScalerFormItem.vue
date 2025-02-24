<template>
  <bcs-form class="form-content" :label-width="210">
    <bcs-form-item
      v-for="item in list"
      :key="item.prop"
      :label="item.name"
      class="form-content-item"
      style="width: 50%;"
      desc-icon="bk-icon icon-info-circle"
      :desc="item.desc">
      <template v-if="item.prop === 'status'">
        <LoadingIcon v-if="autoscalerData[item.prop] === 'UPDATING'">
          {{scalerStatusMap[autoscalerData[item.prop]]}}
        </LoadingIcon>
        <StatusIcon
          :status="autoscalerData[item.prop]"
          :status-color-map="scalerColorMap"
          v-else>
          <span
            :class="{ 'error-tips': autoscalerData[item.prop] === 'UPDATE-FAILURE' }"
            v-bk-tooltips="{
              content: autoscalerData.errorMessage,
              disabled: autoscalerData[item.prop] !== 'UPDATE-FAILURE' || !autoscalerData.errorMessage,
              width: 400
            }">
            {{scalerStatusMap[autoscalerData[item.prop]] || $t('未知')}}
          </span>
        </StatusIcon>
        <template v-if="autoscalerData[item.prop] === 'UPDATE-FAILURE'">
          <span
            class="ml10"
            v-bk-tooltips="$t('查看详情')"
            @click="handleGotoHelmRelease">
            <i class="bcs-icon bcs-icon-fenxiang"></i>
          </span>
        </template>
      </template>
      <span v-else-if="typeof autoscalerData[item.prop] === 'boolean'">
        {{autoscalerData[item.prop] ? $t('是') : $t('否')}}
      </span>
      <span v-else>
        {{`${autoscalerData[item.prop]} ${item.unit || ''}`}}
        <span v-if="item.suffix" class="ml10">{{item.suffix}}</span>
      </span>
      <slot name="suffix" :data="item"></slot>
    </bcs-form-item>
  </bcs-form>
</template>
<script lang="ts">
import { defineComponent, PropType } from '@vue/composition-api';
import $i18n from '@/i18n/i18n-setup';
import StatusIcon from '@/views/dashboard/common/status-icon';
import LoadingIcon from '@/components/loading-icon.vue';
import $router from '@/router/index';
import { BCS_CLUSTER } from '@/common/constant';

export default defineComponent({
  name: 'AutoScalerFormItem',
  components: { StatusIcon, LoadingIcon },
  props: {
    list: {
      type: Array as PropType<any[]>,
      default: () => [],
    },
    autoscalerData: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, ctx) {
    const { $route } = ctx.root;
    // 获取自动扩缩容配置
    const scalerStatusMap = { // 自动扩缩容状态
      NORMAL: $i18n.t('正常'),
      UPDATING: $i18n.t('更新中'),
      'UPDATE-FAILURE': $i18n.t('更新失败'),
      STOPPED: $i18n.t('已停用'),
    };
    const scalerColorMap = {
      NORMAL: 'green',
      UPDATING: 'green',
      'UPDATE-FAILURE': 'red',
      STOPPED: 'gray',
    };
    const handleGotoHelmRelease = () => {
      sessionStorage.setItem(BCS_CLUSTER, $route.params.clusterId);
      $router.push({
        name: 'releaseList',
      });
    };
    return {
      scalerStatusMap,
      scalerColorMap,
      handleGotoHelmRelease,
    };
  },
});
</script>
<style lang="postcss" scoped>
>>> .form-content {
  display: flex;
  flex-wrap: wrap;
  &-item {
      height: 32px;
      margin-top: 0;
      font-size: 12px;
      width: 100%;
  }
  .bk-label {
      font-size: 12px;
      color: #979BA5;
      text-align: left;
  }
  .bk-form-content {
      font-size: 12px;
      color: #313238;
      display: flex;
      align-items: center;
  }
}
>>> .bcs-icon-fenxiang {
    color: #3a84ff;
    cursor: pointer;
}
>>> .error-tips {
    line-height: 1;
    padding: 2px 0;
    border-bottom: 1px dashed #979BA5;
}
</style>
