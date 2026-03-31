import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Modal,
  FlatList,
} from 'react-native';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';

interface DropdownItem {
  id: string;
  label: string;
}

interface DropdownProps {
  items: DropdownItem[];
  selectedId: string | null;
  onSelect: (id: string) => void;
  placeholder?: string;
}

export default function Dropdown({
  items,
  selectedId,
  onSelect,
  placeholder = 'Select...',
}: DropdownProps) {
  const [modalVisible, setModalVisible] = useState(false);

  const selectedItem = items.find(item => item.id === selectedId);

  return (
    <>
      <TouchableOpacity
        style={styles.dropdown}
        onPress={() => setModalVisible(true)}
        activeOpacity={0.7}
      >
        <Text style={[styles.text, !selectedItem && styles.placeholder]}>
          {selectedItem ? selectedItem.label : placeholder}
        </Text>
        <Text style={styles.arrow}>▼</Text>
      </TouchableOpacity>

      <Modal
        animationType="fade"
        transparent
        visible={modalVisible}
        onRequestClose={() => setModalVisible(false)}
      >
        <TouchableOpacity
          style={styles.modalOverlay}
          activeOpacity={1}
          onPress={() => setModalVisible(false)}
        >
          <View style={styles.modalContent}>
            <FlatList
              data={items}
              keyExtractor={item => item.id}
              renderItem={({ item }) => (
                <TouchableOpacity
                  style={styles.modalItem}
                  onPress={() => {
                    onSelect(item.id);
                    setModalVisible(false);
                  }}
                >
                  <Text style={styles.modalItemText}>{item.label}</Text>
                  {item.id === selectedId && (
                    <Text style={styles.checkmark}>✓</Text>
                  )}
                </TouchableOpacity>
              )}
            />
          </View>
        </TouchableOpacity>
      </Modal>
    </>
  );
}

const styles = StyleSheet.create({
  dropdown: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.md,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    ...Shadow.sm,
  },
  text: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
  },
  placeholder: {
    color: Colors.textDisabled,
  },
  arrow: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContent: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    width: '80%',
    maxHeight: '70%',
    padding: Spacing.md,
    ...Shadow.lg,
  },
  modalItem: {
    paddingVertical: Spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  modalItemText: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
  },
  checkmark: {
    color: Colors.primary,
    fontSize: 16,
    fontWeight: 'bold',
  },
});